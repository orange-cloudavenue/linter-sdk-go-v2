package linters

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

var _ register.LinterPlugin = (*PluginSDKV2)(nil)

func init() {
	register.Plugin("sdkv2-api-types", New)
}

type Element struct {
	Name string `json:"name"`
}

type PluginSDKV2 struct {
}

func New(_ any) (register.LinterPlugin, error) {
	// This plugin does not require setting up any specific configuration.
	return &PluginSDKV2{}, nil
}

func (f *PluginSDKV2) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "apitpyesnaming",
			Doc:  "checks that API types follow naming conventions (apiResponse<Object>, apiRequest<Object>, Model<Object>, Params<Object>)",
			Run:  f.runAPITypes,
		},
	}, nil
}

func (f *PluginSDKV2) GetLoadMode() string {
	// NOTE: the mode can be `register.LoadModeSyntax` or `register.LoadModeTypesInfo`.
	// - `register.LoadModeSyntax`: if the linter doesn't use types information.
	// - `register.LoadModeTypesInfo`: if the linter uses types information.

	return register.LoadModeSyntax
}

var (
	apiResponseRe = regexp.MustCompile(`^apiResponse[A-Z][A-Za-z0-9]*$`)
	apiRequestRe  = regexp.MustCompile(`^apiRequest[A-Z][A-Za-z0-9]*$`)
	modelRe       = regexp.MustCompile(`^Model[A-Z][A-Za-z0-9]*$`)
	paramsRe      = regexp.MustCompile(`^Params[A-Z][A-Za-z0-9]*$`)
)

func (f *PluginSDKV2) runAPITypes(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		// Only check files in api/ directory
		if !strings.Contains(pass.Fset.Position(file.Pos()).Filename, "/api/") {
			continue
		}
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				name := typeSpec.Name.Name
				// Only check struct types
				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					continue
				}
				if !apiResponseRe.MatchString(name) &&
					!apiRequestRe.MatchString(name) &&
					!modelRe.MatchString(name) &&
					!paramsRe.MatchString(name) {
					pass.Reportf(typeSpec.Pos(), "type %q does not follow API type naming conventions (see CONTRIBUTING.md)", name)
				}
			}
		}
	}
	return nil, nil
}
