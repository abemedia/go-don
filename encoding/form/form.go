package form

import (
	"net/http"

	"github.com/abemedia/go-don"
	"github.com/gorilla/schema"
)

var MemoryLimit int64 = 1 << 20 // 1MB

var decoder = func() *schema.Decoder {
	d := schema.NewDecoder()
	d.SetAliasTag("form")
	d.IgnoreUnknownKeys(true)
	return d
}()

func decodeForm(r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return decoder.Decode(v, r.Form)
}

func decodeMultipartForm(r *http.Request, v interface{}) error {
	if err := r.ParseMultipartForm(MemoryLimit); err != nil {
		return err
	}
	return decoder.Decode(v, r.Form)
}

func init() {
	don.RegisterDecoder("application/x-www-form-urlencoded", decodeForm)
	don.RegisterDecoder("multipart/form-data", decodeMultipartForm)
}
