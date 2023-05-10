package text

import (
	"github.com/abemedia/go-don/encoding"
)

func init() {
	mediaType := "text/plain"

	encoding.RegisterDecoder(decode, mediaType)
	encoding.RegisterEncoder(encode, mediaType)
}
