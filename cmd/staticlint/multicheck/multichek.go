package multicheck

import (
	"strings"

	"github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/sergalkin/go-url-shortener.git/cmd/staticlint/osexit"
)

type analyzerList []*analysis.Analyzer
type LintOptions func(*analyzerList)

var list analyzerList

// NewWithOptions - creates list of analyzerList based on passed options.
func NewWithOptions(opts ...LintOptions) *analyzerList {
	for _, opt := range opts {
		opt(&list)
	}

	return &list
}

// WithBuiltin - adds standard rules to rules list if -builtin flag was passed.
func WithBuiltin(isActive bool) LintOptions {
	return func(l *analyzerList) {
		if isActive {
			options := []*analysis.Analyzer{
				asmdecl.Analyzer,
				assign.Analyzer,
				atomic.Analyzer,
				atomicalign.Analyzer,
				bools.Analyzer,
				buildssa.Analyzer,
				buildtag.Analyzer,
				cgocall.Analyzer,
				composite.Analyzer,
				copylock.Analyzer,
				ctrlflow.Analyzer,
				deepequalerrors.Analyzer,
				errorsas.Analyzer,
				fieldalignment.Analyzer,
				findcall.Analyzer,
				framepointer.Analyzer,
				httpresponse.Analyzer,
				ifaceassert.Analyzer,
				inspect.Analyzer,
				loopclosure.Analyzer,
				lostcancel.Analyzer,
				nilfunc.Analyzer,
				nilness.Analyzer,
				pkgfact.Analyzer,
				printf.Analyzer,
				reflectvaluecompare.Analyzer,
				shadow.Analyzer,
				shift.Analyzer,
				sigchanyzer.Analyzer,
				sortslice.Analyzer,
				stdmethods.Analyzer,
				stringintconv.Analyzer,
				structtag.Analyzer,
				testinggoroutine.Analyzer,
				tests.Analyzer,
				unmarshal.Analyzer,
				unreachable.Analyzer,
				unsafeptr.Analyzer,
				unusedresult.Analyzer,
				unusedwrite.Analyzer,
				usesgenerics.Analyzer,
				osexit.ExitAnalyzer,
			}
			*l = append(*l, options...)
		}
	}
}

// WithStatic - adds SA rules from https://staticcheck.io to rules list if -static flag was passed.
func WithStatic(isActive bool) LintOptions {
	return func(l *analyzerList) {
		if isActive {
			additionalChecks := map[string]bool{
				"ST1003": true,
				"ST1006": true,
				"ST1023": true,
			}

			for _, v := range staticcheck.Analyzers {
				if strings.Contains(v.Analyzer.Name, "SA") {
					*l = append(*l, v.Analyzer)
				}
			}

			for _, v := range stylecheck.Analyzers {
				if additionalChecks[v.Analyzer.Name] {
					*l = append(*l, v.Analyzer)
				}
			}
		}
	}
}

// WithExtra - adds go-critic and bodyclose to rules list if -extra flag was passed.
func WithExtra(isActive bool) LintOptions {
	return func(l *analyzerList) {
		if isActive {
			*l = append(*l, analyzer.Analyzer)
			*l = append(*l, bodyclose.Analyzer)
		}
	}
}
