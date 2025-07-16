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
