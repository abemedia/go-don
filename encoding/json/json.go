package json

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/goccy/go-json"
)

func init() {
	mediaType := "application/json"

	encoding.RegisterDecoder(json.Unmarshal, mediaType)
	encoding.RegisterEncoder(json.Marshal, mediaType)
}
