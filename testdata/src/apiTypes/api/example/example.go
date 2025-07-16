// nolint: unused
package example

// Valid API types

type (
	apiResponseExample struct{}
	apiRequestExample  struct{}
	ModelExample       struct{}
	ParamsExample      struct{}
	Client             struct{}
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
	ClientExample struct { // want "type ClientExample does not follow API type naming conventions"
	}
	ExampleClient struct { // want "type ExampleClient does not follow API type naming conventions"
	}
)
