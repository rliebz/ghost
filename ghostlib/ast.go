package ghostlib

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"runtime"
	"strings"
	"sync"
)

// Args captures the call site information for deferred AST lookup.
type Args struct {
	once func() []string
}

// Get the string representation of the argument in the specified position.
// Panics on out-of-bounds lookups. Safe for concurrent use.
func (a Args) Get(index int) string {
	return a.once()[index]
}

// ArgsFromAST gets the string representation of the caller's arguments from
// the AST. To handle situations where this cannot be done reliably, the raw
// arguments should be passed so their values can be used as a backup.
//
// This function must be called directly from an assertion function to ensure
// that the intended arguments are captured.
func ArgsFromAST(unformatted ...any) Args {
	pc, _, _, _ := runtime.Caller(1)
	_, file, line, _ := runtime.Caller(2)
	return Args{
		once: sync.OnceValue(func() []string {
			return argsFromAST(pc, findSystemFilepath(file), line, unformatted...)
		}),
	}
}

func argsFromAST(pc uintptr, filename string, line int, unformatted ...any) []string {
	wantFunc := runtime.FuncForPC(pc)
	if wantFunc == nil {
		return mapString(unformatted)
	}

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return mapString(unformatted)
	}

	node := callExprForFunc(wantFunc, fset, astFile, line)
	if node == nil {
		return mapString(unformatted)
	}

	out := make([]string, 0, len(node.Args))
	for _, arg := range node.Args {
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
