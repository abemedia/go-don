// Package msgpack provides encoding and decoding of MessagePack data.
package msgpack

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/vmihailenco/msgpack/v5"
)

func init() {
	mediaType := "application/msgpack"
	aliases := []string{"application/x-msgpack", "application/vnd.msgpack"}

	encoding.RegisterDecoder(msgpack.Unmarshal, mediaType, aliases...)
	encoding.RegisterEncoder(msgpack.Marshal, mediaType, aliases...)
}
