package xml

import (
	"encoding/xml"
	"net/http"

	"github.com/abemedia/go-don"
)

func decodeXML(r *http.Request, v interface{}) error {
	return xml.NewDecoder(r.Body).Decode(v)
}

func encodeXML(w http.ResponseWriter, v interface{}) error {
	return xml.NewEncoder(w).Encode(v)
}

func init() {
	don.RegisterDecoder("application/xml", decodeXML, "text/xml")
	don.RegisterEncoder("application/xml", encodeXML, "text/xml")
}
