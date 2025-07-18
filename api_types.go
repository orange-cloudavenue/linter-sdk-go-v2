package linters

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var (
	// Regex101 https://regex101.com/r/9Mv2Ak
	apiTypesRe = regexp.MustCompile(`(^(apiResponse|apiRequest|Model|Params)[A-Z][A-Za-z0-9]*$|^Client$)`)
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
