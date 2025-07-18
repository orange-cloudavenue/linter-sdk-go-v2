package linters

import (
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
		{
			Name: "endpointstructfields",
			Doc:  "checks that Endpoint struct fields follow naming conventions",
			Run:  f.runEndpointStructFields,
		},
	}, nil
}

func (f *PluginSDKV2) GetLoadMode() string {
	// NOTE: the mode can be `register.LoadModeSyntax` or `register.LoadModeTypesInfo`.
	// - `register.LoadModeSyntax`: if the linter doesn't use types information.
	// - `register.LoadModeTypesInfo`: if the linter uses types information.

	return register.LoadModeSyntax
}
