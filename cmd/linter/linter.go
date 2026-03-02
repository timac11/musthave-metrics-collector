package linter

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "no panic(), no log.Fatal() and no os.exit()",
	Doc:      "no panic(), no log.Fatal() and no os.exit()",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		fn, ok := n.(*ast.FuncDecl)

		if !ok {
			return
		}

		isMain := fn.Name.Name == "main" && pass.Pkg.Name() == "main"

		ast.Inspect(fn, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "panic" {
				pass.Reportf(call.Lparen, "panic call is used")
			}

			if isMain {
				return true
			}

			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				ident, ok := sel.X.(*ast.Ident);

				if !ok {
					return true
				}

				if ident.Name == "os" && sel.Sel.Name == "Exit" {
					pass.Reportf(call.Lparen, "os.Exit is used outside main package")
				}

				if ident.Name == "log" && sel.Sel.Name == "Fatal" {
					pass.Reportf(call.Lparen, "log.Fatal is used outside main package")
				}
			}
			return true
		})
	})
	return nil, nil
}
