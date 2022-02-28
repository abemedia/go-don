package yaml

import (
	"net/http"

	"github.com/abemedia/go-don"
	"gopkg.in/yaml.v2"
)

func decodeYAML(r *http.Request, v interface{}) error {
	return yaml.NewDecoder(r.Body).Decode(v)
}

func encodeYAML(w http.ResponseWriter, v interface{}) error {
	return yaml.NewEncoder(w).Encode(v)
}

func init() {
	don.RegisterEncoder("application/x-yaml", encodeYAML, "text/x-yaml")
	don.RegisterDecoder("application/x-yaml", decodeYAML, "text/x-yaml")
}
