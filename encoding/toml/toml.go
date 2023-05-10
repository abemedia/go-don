package toml

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/pelletier/go-toml"
)

func init() {
	mediaType := "application/toml"

	encoding.RegisterDecoder(toml.Unmarshal, mediaType)
	encoding.RegisterEncoder(toml.Marshal, mediaType)
}
