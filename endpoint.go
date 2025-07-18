package linters

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/orange-cloudavenue/common-go/validators"
	"golang.org/x/tools/go/analysis"
)

func (f *PluginSDKV2) runEndpointStructFields(pass *analysis.Pass) (any, error) {
	// Debug AST traversal https://yuroyoro.github.io/goast-viewer/index.html

	var expectedFields = []string{
		"Name",
		"PathTemplate",
		"Description",
		"Method",
		"DocumentationURL",
	}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.VAR {
				// Match only variable declarations
				continue
			}

			for _, spec := range gen.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				for _, name := range valueSpec.Values {
					compositeLit, ok := name.(*ast.CompositeLit)
					if !ok {
						continue
					}

					typeName, ok := compositeLit.Type.(*ast.Ident)
					if !ok || typeName.Name != "Endpoint" {
						continue
					}

					// ! At this point, we have an Endpoint struct

					// Check if all expected fields are present
					for _, field := range expectedFields {
						found := false
						for _, elt := range compositeLit.Elts {
							kv, ok := elt.(*ast.KeyValueExpr)
							if !ok {
								continue
							}

							fieldName, ok := kv.Key.(*ast.Ident)
							if !ok {
								continue
							}

							if fieldName.Name == field {
								found = true
								break
							}
						}
						if !found {
							pass.Reportf(typeName.Pos(), "field '%s' is missing", field)
						}
					}

					// Store variable for later use
					var (
						pathTemplate  string
						pathParamsPos token.Pos
						pathParams    = make([]struct {
							name string
							pos  token.Pos
						}, 0)
					)

					for _, elt := range compositeLit.Elts {
						kv, ok := elt.(*ast.KeyValueExpr)
						if !ok {
							continue
						}

						fieldName, ok := kv.Key.(*ast.Ident)
						if !ok {
							continue
						}

						// Here we check the field by field
						switch fieldName.Name {
						case "Name":
							fieldValue, err := fieldValueString(kv.Value)
							if err != nil {
								pass.Reportf(kv.Pos(), "field 'Name' must be a string literal")
								continue
							}

							if fieldValue.Value == "" {
								pass.Reportf(fieldValue.Pos(), "field 'Name' cannot be empty")
								continue
							}
							if err := validators.New().Var(fieldValue.Value, "case=PascalCase"); err != nil {
								pass.Reportf(fieldValue.Pos(), "field 'Name' must be in PascalCase")
							}
						case "Method":
							fieldValue, err := fieldValueString(kv.Value)
							if err != nil {
								pass.Reportf(kv.Pos(), "field 'Method' must be a string literal")
								continue
							}

							if fieldValue.Value == "" {
								pass.Reportf(fieldValue.Pos(), "field 'Method' cannot be empty")
								continue
							}

							if err := validators.New().Var(fieldValue.Value, "oneof=GET POST PUT PATCH DELETE"); err != nil {
								pass.Reportf(fieldValue.Pos(), "field 'Method' must be one of GET, POST, PUT, PATCH, DELETE")
							}

						case "PathTemplate":
							fieldValue, err := fieldValueString(kv.Value)
							if err != nil {
								pass.Reportf(kv.Pos(), "field 'Name' must be a string literal")
								continue
							}
							if fieldValue.Value == "" {
								pass.Reportf(fieldValue.Pos(), "field 'PathTemplate' cannot be empty")
								continue
							} else {
								pathTemplate = fieldValue.Value
							}

							if !strings.HasPrefix(pathTemplate, "/") {
								pass.Reportf(fieldValue.Pos(), "field 'PathTemplate' must start with a '/'")
							}
						case "Description":
							fieldValue, err := fieldValueString(kv.Value)
							if err != nil {
								pass.Reportf(kv.Pos(), "field 'Name' must be a string literal")
								continue
							}
							if fieldValue.Value == "" {
								pass.Reportf(fieldValue.Pos(), "field 'Description' cannot be empty")
							}
						case "DocumentationURL":
							fieldValue, err := fieldValueString(kv.Value)
							if err != nil {
								pass.Reportf(kv.Pos(), "field 'Name' must be a string literal")
								continue
							}
							if fieldValue.Value == "" {
								pass.Reportf(fieldValue.Pos(), "field 'DocumentationURL' cannot be empty")
								continue
							}

							if err := validators.New().Var(fieldValue.Value, "http_url"); err != nil {
								pass.Reportf(fieldValue.Pos(), "field 'DocumentationURL' must be a valid URL")
							}
						case "PathParams":
							pathParamsPos = kv.Pos()
							// Check if PathParams is a slice of strings
							pathParamsElts, err := fieldValueSlice(kv.Value, "PathParam")
							if err != nil {
								pass.Reportf(kv.Pos(), "field 'PathParams' must be a slice of PathParam")
								continue
							}

							for _, elt := range pathParamsElts {
								kv, ok := elt.(*ast.KeyValueExpr)
								if !ok {
									continue
								}

								fieldName, ok := kv.Key.(*ast.Ident)
								if !ok {
									continue
								}

								switch fieldName.Name {
								case "Name":

									fieldValue, err := fieldValueString(kv.Value)
									if err != nil {
										pass.Reportf(kv.Pos(), "field 'PathParam.Name' must be a string literal")
										continue
									}
									if fieldValue.Value == "" {
										pass.Reportf(fieldValue.Pos(), "field 'PathParam.Name' cannot be empty")
										continue
									}

									pathParams = append(pathParams, struct {
										name string
										pos  token.Pos
									}{
										name: fieldValue.Value,
										pos:  fieldValue.Pos(),
									})
								case "Description":
									fieldValue, err := fieldValueString(kv.Value)
									if err != nil {
										pass.Reportf(kv.Pos(), "field 'PathParam.Description' must be a string literal")
										continue
									}
									if fieldValue.Value == "" {
										pass.Reportf(fieldValue.Pos(), "field 'PathParam.Description' cannot be empty")
										continue
									}
								}
							}
						case "QueryParams":
							pathParamsPos = kv.Pos()
							// Check if PathParams is a slice of strings
							pathParamsElts, err := fieldValueSlice(kv.Value, "QueryParam")
							if err != nil {
								pass.Reportf(kv.Pos(), "field 'QueryParams' must be a slice of QueryParam")
								continue
							}

							for _, elt := range pathParamsElts {
								kv, ok := elt.(*ast.KeyValueExpr)
								if !ok {
									continue
								}

								fieldName, ok := kv.Key.(*ast.Ident)
								if !ok {
									continue
								}

								switch fieldName.Name {
								case "Name":
									fieldValue, err := fieldValueString(kv.Value)
									if err != nil {
										pass.Reportf(kv.Pos(), "field 'QueryParam.Name' must be a string literal")
										continue
									}
									if fieldValue.Value == "" {
										pass.Reportf(fieldValue.Pos(), "field 'QueryParam.Name' cannot be empty")
										continue
									}

								case "Description":
									fieldValue, err := fieldValueString(kv.Value)
									if err != nil {
										pass.Reportf(kv.Pos(), "field 'QueryParam.Description' must be a string literal")
										continue
									}
									if fieldValue.Value == "" {
										pass.Reportf(fieldValue.Pos(), "field 'QueryParam.Description' cannot be empty")
										continue
									}
								}
							}
						}
					}

					// Check if all path parameters are present in the PathTemplate
					paramInPathTemplate := []string{}
					// Split PathTemplate by '/' and check if each path parameter is present
					for _, segment := range strings.Split(pathTemplate, "/") {
						if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
							paramInPathTemplate = append(paramInPathTemplate, segment[1:len(segment)-1])
						}
					}
					for _, paramName := range paramInPathTemplate {
						found := false
						for _, param := range pathParams {
							if param.name == paramName {
								found = true
								break
							}
						}
						if !found {
							pass.Reportf(pathParamsPos, "PathTemplate contains undeclared path parameter '%s'", paramName)
						}
					}
				}
			}
		}

	}
	return nil, nil
}

func fieldValueString(v ast.Expr) (*ast.BasicLit, error) {
	fieldValue, ok := v.(*ast.BasicLit)
	if !ok || fieldValue.Kind != token.STRING {
		// We only check string literals for field values
		return nil, fmt.Errorf("field value is not a string literal")
	}
	fieldValue.Value = fieldValue.Value[1 : len(fieldValue.Value)-1]

	return fieldValue, nil
}

func fieldValueSlice(v ast.Expr, typeName string) ([]ast.Expr, error) {
	fieldValue, ok := v.(*ast.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("field value is not a slice")
	}

	fieldType, ok := fieldValue.Type.(*ast.ArrayType)
	if !ok {
		return nil, fmt.Errorf("field value is not a slice type")
	}

	fieldTypeElt, ok := fieldType.Elt.(*ast.Ident)
	if !ok {
		return nil, fmt.Errorf("field value slice element type is not an identifier")
	}
	if fieldTypeElt.Name != typeName {
		return nil, fmt.Errorf("field value slice element type is not %s", typeName)
	}

	if len(fieldValue.Elts) == 0 {
		return nil, fmt.Errorf("field value slice is empty")
	}

	return fieldValue.Elts[0].(*ast.CompositeLit).Elts, nil
}
