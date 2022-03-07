package don

import (
	"github.com/gorilla/schema"
)

var (
	decodeHeader = newDecoder("header")
	decodeQuery  = newDecoder("query")
	decodePath   = newDecoder("path")
)

type decoder struct {
	*schema.Decoder
	*tagChecker
}

func newDecoder(tag string) *decoder {
	d := schema.NewDecoder()
	d.SetAliasTag(tag)
	d.IgnoreUnknownKeys(true)

	tc := newTagCheck(tag)

	return &decoder{d, tc}
}
