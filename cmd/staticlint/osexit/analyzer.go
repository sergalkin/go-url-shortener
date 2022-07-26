package osexit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var ExitAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os.Exit in main file",
	Run:  run,
}

type osExitVisitor struct {
	Pass *analysis.Pass
}

func (v *osExitVisitor) Visit(node ast.Node) ast.Visitor {
	if v.Pass.Pkg.Name() != "main" {
		return nil
	}

	// added ignore of some .test file to prevent false positive warnings cause of cache
	if strings.Contains(v.Pass.Pkg.Path(), ".test") {
		return nil
	}

	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Name.Name != "main" {
			return nil
		}
	case *ast.CallExpr:
		v.expr(n)
	}

	return v
}

func (v *osExitVisitor) expr(ce *ast.CallExpr) {
	if fun, funOk := ce.Fun.(*ast.SelectorExpr); funOk {
		pkg, ok := fun.X.(*ast.Ident)
		if !ok || pkg.Name != "os" || fun.Sel.Name != "Exit" {
			return
		}
		v.Pass.Reportf(fun.Pos(), "found os.Exit usage")
	}
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Walk(&osExitVisitor{Pass: pass}, file)
	}
	return nil, nil
}
