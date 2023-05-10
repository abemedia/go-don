package yaml

import (
	"github.com/abemedia/go-don/encoding"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v3"
)

func decodeYAML(ctx *fasthttp.RequestCtx, v any) error {
	return yaml.NewDecoder(ctx.RequestBodyStream()).Decode(v)
}

func encodeYAML(ctx *fasthttp.RequestCtx, v any) error {
	return yaml.NewEncoder(ctx).Encode(v)
}

func init() {
	encoding.RegisterEncoder(encodeYAML, "application/x-yaml", "text/x-yaml")
	encoding.RegisterDecoder(decodeYAML, "application/x-yaml", "text/x-yaml")
}
