# SDKv2 Linters

These linters are designed to enforce specific coding standards and conventions for the CloudAvenue SDK v2. They help maintain code quality and consistency across the SDK.

## What does it check?

### 1. API Types Naming Linter (`apitypesnaming`)

Ensures that struct types in any `api/` directory follow these naming conventions:

- **API Response Types:**  
  Must be named `apiResponse<Object>` (e.g., `apiResponseEdgeGateway`).

- **API Request Body Types:**  
  Must be named `apiRequest<Object>` (e.g., `apiRequestEdgeGateway`).

- **User-facing Model Types:**  
  Must be named `Model<Object>` (e.g., `ModelEdgeGateway`).

- **User-supplied Parameter Types:**  
  Must be named `Params<Object>` (e.g., `ParamsEdgeGateway`).

- **Client Types:**  
  Must be named `Client` (exactly, for the main client struct of an API group, e.g., `type Client struct { ... }`).

If a struct type in an `api/` directory does not follow one of these conventions, the linter will report an error.

#### Example

```go
type EdgeGatewayResponse struct { // ❌ This will trigger a linter warning!
    // ...
}

type apiResponseEdgeGateway struct { // ✅ This is correct!
    // ...
}

type Client struct { // ✅ This is correct for the main API client!
    // ...
}
```

#### Notes

- The `Client` type should be named exactly `Client` (not `ClientEdgeGateway` or similar) for the main client struct in each API group.
- All other types should follow the `apiResponse<Object>`, `apiRequest<Object>`, `Model<Object>`, or `Params<Object>` naming conventions as appropriate.

#### Regex Rule

You can visualize and debug the naming convention regex used by this linter at the following link:  
[https://regex101.com/r/g8Av6t](https://regex101.com/r/g8Av6t)

This tool helps you understand how the linter matches valid type names and can assist in troubleshooting naming issues.

---

### 2. Endpoint Struct Fields Linter (`endpointstructfields`)

Validates the usage and values of the `Endpoint` struct in the `cav/` package.  
It checks for:

- Presence of required fields: `Name`, `PathTemplate`, `Description`, `Method`, `DocumentationURL`
- Field value correctness:
  - `Name` must be non-empty and in PascalCase
  - `Description` must be non-empty
  - `Method` must be one of `GET`, `POST`, `PUT`, `PATCH`, `DELETE`
  - `PathTemplate` must be non-empty and start with `/`
  - `DocumentationURL` must be a valid URL
- Path and query parameters:
  - Each `PathParam` and `QueryParam` must have non-empty `Name` and `Description`
  - All path parameters in `PathTemplate` must be declared in `PathParams`

#### Example

```go
_ = Endpoint{
    Name:             "MyEndpoint", // ✅ PascalCase
    Description:      "This is an example endpoint",
    Method:           "GET",        // ✅ Allowed method
    PathTemplate:     "/v1/example/{id}", // ✅ Starts with /
    DocumentationURL: "https://docs.example.com/api/v1/example", // ✅ Valid URL
    PathParams: []PathParam{
        {
            Name:        "id", // ✅ Declared and used in PathTemplate
            Description: "The ID of the resource",
        },
    },
}
```

If any required field is missing or invalid, or if a path parameter is undeclared, the linter will report an error.

---

## How to Use

These linters are implemented as a golangci-lint plugin.  
Register them in your `.golangci.yml` and ensure they are enabled for your project.

## References

- See [CONTRIBUTING.md](../CONTRIBUTING.md) for the full naming convention rules.
- For more on custom linters, see the [golangci-lint documentation](https://golangci-lint.run/usage/plugins/).
