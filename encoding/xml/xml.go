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
	encoding.RegisterDecoder(decodeXML, "application/xml", "text/xml")
	encoding.RegisterEncoder(encodeXML, "application/xml", "text/xml")
}
