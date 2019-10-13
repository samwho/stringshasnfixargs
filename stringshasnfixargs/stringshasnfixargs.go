package stringshasnfixargs

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "stringshasnfixargs",
	Doc:      "reports potentially incorrect usages of strings.Has{Prefix,Suffix}",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)

		if !isRelevant(pass, call) {
			return
		}

		s := pass.TypesInfo.Types[call.Args[0]]
		prefix := pass.TypesInfo.Types[call.Args[1]]

		if prefix.Value == nil && s.Value != nil {
			pass.Reportf(call.Pos(), "found constant s and variable prefix, arguments to this function may be the wrong way around")
		}
	})

	return nil, nil
}

func isRelevant(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if sel.Sel.Name != "HasPrefix" && sel.Sel.Name != "HasSuffix" {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	if ident.Name != "strings" {
		return false
	}

	return true
}
