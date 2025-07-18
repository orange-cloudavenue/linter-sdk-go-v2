package endpoint

type (
	API     string
	Version string
	Method  string

	Endpoint struct {
		// ID is the unique identifier of the endpoint.
		// It is used to uniquely identify the endpoint in the registry.
		ID string

		// api is the api of the endpoint, e.g., "vdc", "edgegateway", "vapp"
		api API `validate:"required"`

		// version is the API version, e.g., "v1", "v2, "v3", etc.
		// It is used to differentiate between different versions of the API.
		version Version `validate:"required"`

		// Name is the name of the endpoint, e.g., "firewall", "loadBalancer", etc.
		// It is used to group endpoints by their functionality.
		// For example, all endpoints related to firewall operations can be grouped under the "firewall" name.
		Name string `validate:"required"`

		// Description is a brief description of the endpoint.
		// It provides additional information about the endpoint's purpose and functionality.
		// Description is used to provide context in the error messages.
		Description string `validate:"required"`

		// SubClient is the name of the sub-client that this endpoint belongs to.
		// SubClient SubClientName `validate:"required"`

		// Method is the HTTP method used for the endpoint, e.g., "GET", "POST", "PUT", "DELETE".
		Method Method `validate:"required,oneof=GET POST PUT DELETE PATCH"`

		// PathTemplate is the URL path template for the endpoint.
		PathTemplate string `validate:"required"` // e.g., "/v1/edgeGateways/{gatewayId}/firewall/rules"

		// PathParams List of path parameters that can be used in the URL path.
		// These parameters are placeholders in the URL that can be replaced with actual values.
		// For example, in the URL "/v1/edgeGateways/{gatewayId}/firewall/rules",
		// "{gatewayId}" is a path parameter that can be replaced with an actual gateway ID.
		// PathParams are used to construct the final URL for the endpoint.
		PathParams []PathParam `validate:"dive"`

		// QueryParams List of query parameters that can be used in the URL.
		// These parameters are appended to the URL as key-value pairs.
		// For example, in the URL "/v1/edgeGateways/{gatewayId}/firewall/rules?active=true",
		// "active" is a query parameter that can be used to filter results.
		// QueryParams are used to add additional information to the URL for the endpoint.
		QueryParams []QueryParam `validate:"dive"`

		// DocumentationURL is the URL to the documentation for this endpoint.
		DocumentationURL string `validate:"required,url"` // e.g., "https://docs.xx.com/api/v1/xx"

		// BodyRequestType is the golang type of the request body.
		// It is used to validate the body arguments passed to the endpoint.
		// BodyType is optional and can be used to specify the type of the request body
		// for POST, PUT, or PATCH requests.
		BodyRequestType any `validate:"-"`

		// BodyResponseType is the golang type of the response body.
		// It is used to validate the response body returned by the endpoint.
		//
		// If your set `cav.Job{}` as BodyResponseType, the system will automatically
		// handle the job response and retrieve the job status until it is completed (success or error).
		BodyResponseType any `validate:"-"`

		// * Job

		// jobOptions is the options for the job.
		// It is used to specify the options for the job, such as the Timeout, PollingInterval, and ExtractorFunc
		// JobOptions *JobOptions
	}

	QueryParam struct {
		Name          string `validate:"required"`
		Description   string `validate:"required"`
		Required      bool
		ValidatorFunc func(value string) error

		// Ability to provides a value for the query parameter.
		// This is useful when the query parameter value is known at the time of registration.
		// For example, if the query parameter is {type}, you can provide a value like "firewall" or "loadBalancer".
		// If the query parameter is provided Required and ValidatorFunc are ignored.
		Value string
	}

	PathParam struct {
		Name          string `validate:"required"`
		Description   string `validate:"required"`
		Required      bool
		ValidatorFunc func(value string) error

		// Ability to provides a value for the path parameter.
		// This is useful when the path parameter value is known at the time of registration.
		// For example, if the path parameter is {type}, you can provide a value like "firewall" or "loadBalancer".
		// If the path parameter is provided Required and ValidatorFunc are ignored.
		Value string
	}
)

var (
	_ = Endpoint{ // want "field 'Name' is missing" "field 'Description' is missing" "field 'Method' is missing" "field 'PathTemplate' is missing" "field 'DocumentationURL' is missing"
	}

	_ = Endpoint{
		Name:             "", // want "field 'Name' cannot be empty"
		Description:      "", // want "field 'Description' cannot be empty"
		Method:           "", // want "field 'Method' cannot be empty"
		PathTemplate:     "", // want "field 'PathTemplate' cannot be empty"
		DocumentationURL: "", // want "field 'DocumentationURL' cannot be empty"
		PathParams: []PathParam{
			{
				Name:        "", // want "field 'PathParam.Name' cannot be empty"
				Description: "", // want "field 'PathParam.Description' cannot be empty"
			},
		},
		QueryParams: []QueryParam{
			{
				Name:        "", // want "field 'QueryParam.Name' cannot be empty"
				Description: "", // want "field 'QueryParam.Description' cannot be empty"
			},
		},
	}

	_ = Endpoint{
		Name:             "exampleEndpoint", // want "field 'Name' must be in PascalCase"
		Description:      "This is an example endpoint",
		Method:           "HEAD",                            // want "field 'Method' must be one of GET, POST, PUT, PATCH, DELETE"
		PathTemplate:     "v1/example/{id}/{id2}",           // want "field 'PathTemplate' must start with a '/'"
		DocumentationURL: "docs.example.com/api/v1/example", // want "field 'DocumentationURL' must be a valid URL"
		PathParams: []PathParam{ // want "PathTemplate contains undeclared path parameter 'id2'"
			{
				Name:        "id",
				Description: "The ID of the resource",
				Required:    true,
			},
		},
	}
)
