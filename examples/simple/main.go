package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/json"
)

func Ping(ctx context.Context, _ struct{}) (interface{}, error) {
	return "pong", nil
}

type NameRequest struct {
	Name string `json:"name"`         // Get name from JSON body.
	Age  int    `header:"X-User-Age"` // Get age from HTTP header.
}

type NameResponse struct {
	Greeting string `json:"data"`
}

// Set a custom HTTP response code.
func (nr *NameResponse) StatusCode() int {
	return http.StatusTeapot
}

// Add custom headers to the response.
func (nr *NameResponse) Header() http.Header {
	header := http.Header{}
	header.Set("foo", "bar")
	return header
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
	r.Get("/", don.H(Ping))

	g := r.Group("/api")
	g.Get("/ping", don.H(Ping))
	g.Post("/greet", don.H(Greet))

	http.ListenAndServe(":8080", r.Router())
}
