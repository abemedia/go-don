package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/json"
	_ "github.com/abemedia/go-don/encoding/text"
	_ "github.com/abemedia/go-don/encoding/yaml"
)

// Empty returns 204 - No Content.
func Empty(context.Context, any) (any, error) {
	return nil, nil
}

func Ping(context.Context, any) (string, error) {
	return "pong", nil
}

type GreetRequest struct {
	Name string `json:"name"`         // Get name from JSON body.
	Age  int    `header:"X-User-Age"` // Get age from HTTP header.
}

type GreetResponse struct {
	Greeting string `json:"data"`
}

// StatusCode sets a custom HTTP response code.
func (gr *GreetResponse) StatusCode() int {
	return http.StatusTeapot
}

// Header adds custom headers to the response.
func (gr *GreetResponse) Header() http.Header {
	header := http.Header{}
	header.Set("Foo", "bar")
	return header
}

func Greet(_ context.Context, req *GreetRequest) (*GreetResponse, error) {
	if req.Name == "" {
		return nil, don.ErrBadRequest
	}

	res := &GreetResponse{
		Greeting: fmt.Sprintf("Hello %s, you're %d years old.", req.Name, req.Age),
	}

	return res, nil
}

func main() {
	r := don.New(nil)
	r.Get("/", don.H(Empty))

	g := r.Group("/api")
	g.Get("/ping", don.H(Ping))
	g.Post("/greet", don.H(Greet))

	log.Fatal(r.ListenAndServe(":8080"))
}
