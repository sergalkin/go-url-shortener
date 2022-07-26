// Package osexit - is custom analysis rule that searches for os.Exit calls in main function of package main.
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

// Visit - method is invoked for each node encountered by ast.Walk.
func (v *osExitVisitor) Visit(node ast.Node) ast.Visitor {
	if v.Pass.Pkg.Name() != "main" {
		return nil
	}

	// added ignore of some .test file to prevent false positive warnings cause of cache.
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

// expr - is called then node represents an expression
// expr - checks for expression to be os.Exit call with v.Pass.Reportf() call if so.
func (v *osExitVisitor) expr(ce *ast.CallExpr) {
	if fun, funOk := ce.Fun.(*ast.SelectorExpr); funOk {
		pkg, ok := fun.X.(*ast.Ident)
		if !ok || pkg.Name != "os" || fun.Sel.Name != "Exit" {
			return
		}
		v.Pass.Reportf(fun.Pos(), "found os.Exit usage")
	}
}

// run - main function of analyzer.
func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Walk(&osExitVisitor{Pass: pass}, file)
	}
	return nil, nil
}
