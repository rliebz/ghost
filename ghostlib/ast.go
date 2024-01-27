package ghostlib

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"runtime"
	"strings"
)

// ArgsFromAST gets the string representation of the caller's arguments from
// the AST. To handle situations where this cannot be done reliably, the raw
// arguments should be passed so their values can be used as a backup.
func ArgsFromAST(unformatted ...any) []string {
	return argsFromASTSkip(1, unformatted...)
}

// argsFromASTSkip gets the string representation of the caller's arguments
// from the AST, skipping the number specified.
func argsFromASTSkip(skip int, unformatted ...any) []string {
	args, err := callExprArgs(2 + skip)
	if err != nil {
		return mapString(unformatted)
	}

	out := make([]string, 0, len(args))
	for _, arg := range args {
		out = append(out, nodeToString(arg))
	}

	return out
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
		return nil, errors.New("failed to get file/line")
	}

	_, filename, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil, errors.New("failed to get file/line")
	}

	filename = findSystemFilepath(filename)

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

// Passing the -trimpath flag will prevent looking up filepaths directly.
// In most cases, some suffix of the path will be a valid relative path, which
// we can use instead.
func findSystemFilepath(filename string) string {
	if _, err := os.Stat(filename); err == nil {
		return filename
	}

	cur := filename
	for {
		parts := strings.SplitN(cur, "/", 2)
		if len(parts) < 2 {
			break
		}

		cur = parts[1]
		if _, err := os.Stat(cur); err == nil {
			return cur
		}
	}

	return filename
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

func nodeToString(node ast.Node) string {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), node); err != nil {
		panic(err)
	}
	return buf.String()
}
