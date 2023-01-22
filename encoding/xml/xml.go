package xml

import (
	"encoding/xml"

	"github.com/abemedia/go-don"
	"github.com/valyala/fasthttp"
)

func decodeXML(ctx *fasthttp.RequestCtx, v any) error {
	return xml.NewDecoder(ctx.RequestBodyStream()).Decode(v)
}

func encodeXML(ctx *fasthttp.RequestCtx, v any) error {
	return xml.NewEncoder(ctx).Encode(v)
}

func init() {
	don.RegisterDecoder("application/xml", decodeXML, "text/xml")
	don.RegisterEncoder("application/xml", encodeXML, "text/xml")
}
