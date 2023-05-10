package json

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/goccy/go-json"
)

func init() {
	encoding.RegisterDecoder(json.Unmarshal, "application/json")
	encoding.RegisterEncoder(json.Marshal, "application/json")
}
