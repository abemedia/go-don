package yaml

import (
	"github.com/abemedia/go-don"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v3"
)

func decodeYAML(ctx *fasthttp.RequestCtx, v interface{}) error {
	return yaml.NewDecoder(ctx.RequestBodyStream()).Decode(v)
}

func encodeYAML(ctx *fasthttp.RequestCtx, v interface{}) error {
	return yaml.NewEncoder(ctx).Encode(v)
}

func init() {
	don.RegisterEncoder("application/x-yaml", encodeYAML, "text/x-yaml")
	don.RegisterDecoder("application/x-yaml", decodeYAML, "text/x-yaml")
}
