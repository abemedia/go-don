package text

import (
	"encoding"
	"fmt"
	"strconv"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/internal/byteconv"
	"github.com/valyala/fasthttp"
)

//nolint:cyclop
func encode(ctx *fasthttp.RequestCtx, v any) error {
	var (
		b   []byte
		err error
	)

	if v != nil {
		switch v := v.(type) {
		case string:
			b = byteconv.Atob(v)
		case []byte:
			b = v
		case int:
			b = strconv.AppendInt(ctx.Response.Body(), int64(v), 10)
		case int8:
			b = strconv.AppendInt(ctx.Response.Body(), int64(v), 10)
		case int16:
			b = strconv.AppendInt(ctx.Response.Body(), int64(v), 10)
		case int32:
			b = strconv.AppendInt(ctx.Response.Body(), int64(v), 10)
		case int64:
			b = strconv.AppendInt(ctx.Response.Body(), v, 10)
		case uint:
			b = strconv.AppendUint(ctx.Response.Body(), uint64(v), 10)
		case uint8:
			b = strconv.AppendUint(ctx.Response.Body(), uint64(v), 10)
		case uint16:
			b = strconv.AppendUint(ctx.Response.Body(), uint64(v), 10)
		case uint32:
			b = strconv.AppendUint(ctx.Response.Body(), uint64(v), 10)
		case uint64:
			b = strconv.AppendUint(ctx.Response.Body(), v, 10)
		case float32:
			b = strconv.AppendFloat(ctx.Response.Body(), float64(v), 'f', -1, 32)
		case float64:
			b = strconv.AppendFloat(ctx.Response.Body(), v, 'f', -1, 64)
		case bool:
			b = strconv.AppendBool(ctx.Response.Body(), v)
		case error:
			b = byteconv.Atob(v.Error())
		case encoding.TextMarshaler:
			b, err = v.MarshalText()
		case fmt.Stringer:
			b = append(ctx.Response.Body(), v.String()...)
		default:
			return don.ErrNotAcceptable
		}
	}

	if len(b) > 0 {
		ctx.Response.SetBodyRaw(append(b, '\n'))
	}

	return err
}

func init() {
	don.RegisterEncoder("text/plain", encode)
}
