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
	mediaType := "application/yaml"
	aliases := []string{"text/yaml", "application/x-yaml", "text/x-yaml", "text/vnd.yaml"}

	encoding.RegisterDecoder(decodeYAML, mediaType, aliases...)
	encoding.RegisterEncoder(encodeYAML, mediaType, aliases...)
}
