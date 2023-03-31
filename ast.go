package ghost

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"runtime"
	"strings"
)

func getFormattedArgs(skip int) ([]string, bool) {
	args, err := callExprArgs(skip + 1)
	if err != nil {
		return nil, false
	}

	out := make([]string, 0, len(args))
	for _, arg := range args {
		out = append(out, nodeToString(arg))
	}

	return out, true
}

func callExprArgs(skip int) ([]ast.Expr, error) {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return nil, errors.New("failed to get file/line")
	}

	_, filename, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil, errors.New("failed to get file/line")
	}

	wantFunc := runtime.FuncForPC(pc)

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	node := callExprForFunc(wantFunc, fset, astFile, line)
	if node == nil {
		return nil, errors.New("no node found at line")
	}

	return node.Args, nil
}

func callExprForFunc(
	wantFunc *runtime.Func,
	fset *token.FileSet,
	file *ast.File,
	lineNum int,
) *ast.CallExpr {
	var out *ast.CallExpr
	ast.Inspect(file, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		if fset.Position(node.Pos()).Line != lineNum {
			return true
		}

		callExpr, ok := node.(*ast.CallExpr)
		if ok && describesCallExpr(wantFunc, callExpr, fset, file) {
			out = callExpr
		}

		return true
	})
	return out
}

// This comparison isn't perfect, but it works well enough so far.
func describesCallExpr(
	wantFn *runtime.Func,
	callExpr *ast.CallExpr,
	fset *token.FileSet,
	file *ast.File,
) bool {
	wantName := wantFn.Name()
	wantName = strings.TrimSuffix(wantName, "[...]")
	dotIndex := strings.LastIndex(wantName, ".")
	wantName = wantName[dotIndex+1:]

	switch fun := callExpr.Fun.(type) {
	case *ast.Ident:
		return wantName == fun.Name
	case *ast.SelectorExpr:
		return wantName == fun.Sel.Name
	}

	return false
}

func nodeToString(node ast.Node) string {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), node); err != nil {
		panic(err)
	}
	return buf.String()
}
