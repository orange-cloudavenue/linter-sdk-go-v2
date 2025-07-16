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
	register.Plugin("linter-sdkv2", New)
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
			Name: "apitypesnaming",
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
	apiTypesRe = regexp.MustCompile(`(apiResponse|apiRequest|Model|Params|Client)[A-Z][A-Za-z0-9]*$`)
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

				if !apiTypesRe.MatchString(name) {
					pass.Reportf(typeSpec.Pos(), "type %s does not follow API type naming conventions", name)
				}
			}
		}
	}
	return nil, nil
}
