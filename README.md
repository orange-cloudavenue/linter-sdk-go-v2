# SDKv2 Linters

These linters are designed to enforce specific coding standards and conventions for the CloudAvenue SDK v2. They help maintain code quality and consistency across the SDK.

## What does it check?

The linter ensures that struct types in any `api/` directory follow these naming conventions:

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

## Example

Suppose you have the following type in `api/edgegateway/v1/edgegateway.go`:

```go
type EdgeGatewayResponse struct { // ❌ This will trigger a linter warning!
    // ...
}
```

You should rename it to:

```go
type apiResponseEdgeGateway struct { // ✅ This is correct!
    // ...
}
```

## Notes

- The `Client` type should be named exactly `Client` (not `ClientEdgeGateway` or similar) for the main client struct in each API group.
- All other types should follow the `apiResponse<Object>`, `apiRequest<Object>`, `Model<Object>`, or `Params<Object>` naming conventions as appropriate.

## Regex Rule

You can visualize and debug the naming convention regex used by this linter at the following link:  
[https://regex101.com/r/g8Av6t](https://regex101.com/r/g8Av6t)

This tool helps you understand how the linter matches valid type names and can assist in troubleshooting naming issues.
