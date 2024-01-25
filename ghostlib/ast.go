package ghostlib

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"runtime"
	"strings"
)

// ArgsFromAST gets the string representation of the caller's arguments from
// the AST. To handle situations where this cannot be done reliably, the raw
// arguments should be passed so their values can be used as a backup.
func ArgsFromAST(unformatted ...any) []string {
	args, err := argsFromASTSkip(1)
	if err != nil {
		_ = mapString(unformatted) // return
		panic(fmt.Sprintf("failed to parse args from the AST: %s", err))
	}

	return args
}

// argsFromASTSkip gets the string representation of the caller's arguments
// from the AST, skipping the number specified.
func argsFromASTSkip(skip int) ([]string, error) {
	args, err := callExprArgs(2 + skip)
	if err != nil {
		return nil, fmt.Errorf("getting call expr args: %w", err)
	}

	out := make([]string, 0, len(args))
	for i, arg := range args {
		s, err := nodeToString(arg)
		if err != nil {
			return nil, fmt.Errorf("converting node %d to string: %w", i, err)
		}

		out = append(out, s)
	}

	return out, nil
}

func mapString(s []any) []string {
	out := make([]string, 0, len(s))
	for _, ss := range s {
		out = append(out, fmt.Sprint(ss))
	}
	return out
}

func callExprArgs(skip int) ([]ast.Expr, error) {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return nil, errors.New("failed to get file/line for skip")
	}

	_, filename, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil, errors.New("failed to get file/line for skip + 1")
	}

	wantFunc := runtime.FuncForPC(pc)

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
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
		if ok && describesCallExpr(wantFunc, callExpr) {
			out = callExpr
		}

		return true
	})
	return out
}

// This comparison isn't perfect, but it works well enough so far.
func describesCallExpr(wantFn *runtime.Func, callExpr *ast.CallExpr) bool {
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

func nodeToString(node ast.Node) (string, error) {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), node); err != nil {
		return "", fmt.Errorf("formatting node: %w", err)
	}
	return buf.String(), nil
}
