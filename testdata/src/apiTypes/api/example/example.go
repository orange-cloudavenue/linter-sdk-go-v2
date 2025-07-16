// nolint: unused
package example

// Valid API types

type (
	apiResponseExample struct{}
	apiRequestExample  struct{}
	ModelExample       struct{}
	ParamsExample      struct{}
	ClientExample      struct{}
)

// Invalid API types

type (
	exampleResponse struct { // want "type exampleResponse does not follow API type naming conventions"
	}
	exampleRequest struct { // want "type exampleRequest does not follow API type naming conventions"
	}
	exampleModel struct { // want "type exampleModel does not follow API type naming conventions"
	}
	exampleParams struct { // want "type exampleParams does not follow API type naming conventions"
	}
	exampleClient struct { // want "type exampleClient does not follow API type naming conventions"
	}
)
