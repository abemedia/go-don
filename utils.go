package don

import (
	"net/http"
	"strings"
)

func parseMime(contentType, fallback string) string {
	index := strings.Index(contentType, ";")
	if index > 0 {
		contentType = contentType[:index]
	}

	if contentType == "" || contentType == "*/*" {
		return fallback
	}

	return strings.TrimSpace(contentType)
}

func withConfig(handle http.Handler, c *Config) http.Handler {
	if h, ok := handle.(interface{ setConfig(*Config) }); ok {
		h.setConfig(c)
	}
	return handle
}
