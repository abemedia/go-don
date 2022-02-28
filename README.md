# Don - Go API Framework

[![GoDoc](https://pkg.go.dev/badge/github.com/abemedia/go-don)](https://pkg.go.dev/github.com/abemedia/go-don)

Don is a blazing-fast API framework written in Go. It uses the new Go generics and requires Go 1.18 to work.

It's still very early alpha and is likely to change so not recommended for production yet.

## Basic Example

```go
package main

import (
  "context"
  "fmt"
  "net/http"

  "github.com/abemedia/go-don"
  _ "github.com/abemedia/go-don/encoding/form" // Enable form parsing.
  _ "github.com/abemedia/go-don/encoding/json" // Enable JSON parsing & rendering.
  _ "github.com/abemedia/go-don/encoding/yaml" // Enable YAML parsing & rendering.
)

type NameRequest struct {
  Name string `path:"name"`         // Get name from the URL path.
  Age  int    `header:"X-User-Age"` // Get age from HTTP header.
}

type NameResponse struct {
  Greeting string `json:"data" yaml:"data"` // Remember to add all the tags for the renderers you enable.
}

func Greet(ctx context.Context, request NameRequest) (interface{}, error) {
  if request.Name == "" {
    return nil, don.ErrBadRequest
  }

  res := &NameResponse{
    Greeting: fmt.Sprintf("Hello %s, you're %d years old.", request.Name, request.Age),
  }

  return res, nil
}

func main() {
  r := don.New(nil)
  r.Post("/greet/:name", don.H(Greet)) // Handlers are wrapped with `don.H`.
  http.ListenAndServe(":8080", r.Router())
}

```

## Support multiple formats

Support multiple formats without writing extra rendering or parsing code.
The API uses `Content-Type` and `Accept` headers to determine what input and output encoding to use.

Currently supported formats:

- JSON
- XML
- YAML
- Form (input only - both `application/x-www-form-urlencoded` and `multipart/form-data`)

Formats need to be explicitly imported e.g.

```go
import _ "github.com/abemedia/go-don/encoding/yaml"
```

Adding your own is easy. See [encoding/json/json.go](./blob/master/encoding/json/json.go).

## Simple request parsing

Automatically unmarshals values from headers, URL query, URL path & request body into your request struct.

```go
type MyRequest struct {
  // Get from the URL path.
  ID int64 `path:"id"`                                     

  // Get from the URL query.
  Filter string `query:"filter"`                                

  // Get from the JSON, YAML, XML or form body.
  Content float64 `form:"bar" json:"bar" yaml:"bar" xml:"bar"`

  // Get from the HTTP header.
  Lang string `header:"Accept-Language"`                   
}
```

## Customising headers & response codes

Implement the `StatusCoder` and `Headerer` interfaces to customise headers and response codes.

```go
type MyResponse struct {
  Foo  string `json:"foo"`
}

// Set a custom HTTP response code.
func (nr *MyResponse) StatusCode() int {
  return 201
}

// Add custom headers to the response.
func (nr *MyResponse) Header() http.Header {
  header := http.Header{}
  header.Set("foo", "bar")
  return header
}
```

## Middleware

Don uses the standard library middleware format of `func(http.Handler) http.Handler`.

For example:

```go
func loggingMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Println(r.RequestURI)
    next.ServeHTTP(w, r)
  })
}
```

It is registered using `don.Use` e.g.

```go
r := don.New(nil)
r.Post("/", don.H(handler))
r.Use(loggingMiddleware)
```

To pass values from the middleware to the handler extend the context e.g.

```go
func myMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    ctx := context.WithValue(r.Context(), ContextUserKey, "my_user")
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}
```

This can now be accessed in the handler:

```go
user := ctx.Value(ContextUserKey).(string)
```
