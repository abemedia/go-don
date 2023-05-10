package msgpack

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/vmihailenco/msgpack/v5"
)

func init() {
	encoding.RegisterDecoder(msgpack.Unmarshal, "application/x-msgpack")
	encoding.RegisterEncoder(msgpack.Marshal, "application/x-msgpack")
}
