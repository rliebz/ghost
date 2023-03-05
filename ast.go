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
	pc, _, _, _ := runtime.Caller(skip)
	_, filename, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil, errors.New("failed to get file/line")
	}

	fn := runtime.FuncForPC(pc)

	fileset := token.NewFileSet()
	astFile, err := parser.ParseFile(fileset, filename, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	node := callExprForFunc(fn, fileset, astFile, line)
	if node == nil {
		return nil, errors.New("no node found at line")
	}

	return node.Args, nil
}

func callExprForFunc(
	want *runtime.Func,
	fileset *token.FileSet,
	node ast.Node,
	lineNum int,
) *ast.CallExpr {
	var out *ast.CallExpr
	ast.Inspect(node, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		if fileset.Position(node.Pos()).Line != lineNum {
			return true
		}

		if callExpr, ok := node.(*ast.CallExpr); ok && describesCallExpr(want, callExpr) {
			out = callExpr
		}

		return true
	})
	return out
}

// TODO: There are likely more reliable ways to compare functions.
func describesCallExpr(want *runtime.Func, callExpr *ast.CallExpr) bool {
	switch fun := callExpr.Fun.(type) {
	case *ast.Ident:
		return strings.Contains(want.Name(), fun.Name)
	case *ast.SelectorExpr:
		return strings.Contains(want.Name(), nodeToString(fun))
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
