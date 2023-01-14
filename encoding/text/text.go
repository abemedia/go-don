package json

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/abemedia/go-don"
	"github.com/valyala/fasthttp"
)

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//nolint:cyclop
func decodeText(ctx *fasthttp.RequestCtx, v interface{}) error {
	b := ctx.Request.Body()
	if len(b) == 0 {
		return nil
	}

	switch t := v.(type) {
	case *string:
		*t = b2s(b)

	case *[]byte:
		*t = b

	case *int:
		d, err := strconv.Atoi(b2s(b))
		if err != nil {
			return err
		}

		*t = d

	case *int8:
		d, err := strconv.ParseInt(b2s(b), 10, 8)
		if err != nil {
			return err
		}

		*t = int8(d)

	case *int16:
		d, err := strconv.ParseInt(b2s(b), 10, 16)
		if err != nil {
			return err
		}

		*t = int16(d)

	case *int32:
		d, err := strconv.ParseInt(b2s(b), 10, 32)
		if err != nil {
			return err
		}

		*t = int32(d)

	case *int64:
		d, err := strconv.ParseInt(b2s(b), 10, 64)
		if err != nil {
			return err
		}

		*t = d

	case *uint:
		d, err := strconv.ParseUint(b2s(b), 10, 0)
		if err != nil {
			return err
		}

		*t = uint(d)

	case *uint8:
		d, err := strconv.ParseUint(b2s(b), 10, 8)
		if err != nil {
			return err
		}

		*t = uint8(d)

	case *uint16:
		d, err := strconv.ParseUint(b2s(b), 10, 16)
		if err != nil {
			return err
		}

		*t = uint16(d)

	case *uint32:
		d, err := strconv.ParseUint(b2s(b), 10, 32)
		if err != nil {
			return err
		}

		*t = uint32(d)

	case *uint64:
		d, err := strconv.ParseUint(b2s(b), 10, 64)
		if err != nil {
			return err
		}

		*t = d

	case *float32:
		d, err := strconv.ParseFloat(b2s(b), 32)
		if err != nil {
			return err
		}

		*t = float32(d)

	case *float64:
		d, err := strconv.ParseFloat(b2s(b), 64)
		if err != nil {
			return err
		}

		*t = d

	case *bool:
		d, err := strconv.ParseBool(b2s(b))
		if err != nil {
			return err
		}

		*t = d

	default:
		return don.ErrUnsupportedMediaType
	}

	return nil
}

func encodeText(ctx *fasthttp.RequestCtx, v interface{}) error {
	if v != nil {
		switch v.(type) {
		case *string, string,
			*[]byte, []byte,
			*int, int,
			*int8, int8,
			*int16, int16,
			*int32, int32,
			*int64, int64,
			*uint, uint,
			*uint8, uint8,
			*uint16, uint16,
			*uint32, uint32,
			*uint64, uint64,
			*float32, float32,
			*float64, float64,
			*bool, bool,
			error:

		default:
			return don.ErrNotAcceptable
		}
	}

	_, err := fmt.Fprintln(ctx, v)

	return err
}

func init() {
	don.RegisterDecoder("text/plain", decodeText)
	don.RegisterEncoder("text/plain", encodeText)
}
