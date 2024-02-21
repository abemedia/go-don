// Package xml provides encoding and decoding of XML data.
package xml

import (
	"encoding/xml"

	"github.com/abemedia/go-don/encoding"
	"github.com/valyala/fasthttp"
)

func decodeXML(ctx *fasthttp.RequestCtx, v any) error {
	return xml.NewDecoder(ctx.RequestBodyStream()).Decode(v)
}

func encodeXML(ctx *fasthttp.RequestCtx, v any) error {
	return xml.NewEncoder(ctx).Encode(v)
}

func init() {
	mediaType := "application/xml"
	alias := "text/xml"

	encoding.RegisterDecoder(decodeXML, mediaType, alias)
	encoding.RegisterEncoder(encodeXML, mediaType, alias)
}
