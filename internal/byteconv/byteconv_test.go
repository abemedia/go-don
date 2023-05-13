package byteconv_test

import (
	"bytes"
	"testing"

	"github.com/abemedia/go-don/internal/byteconv"
)

func TestBtoa(t *testing.T) {
	if byteconv.Btoa([]byte("test")) != "test" {
		t.Error("should be equal")
	}
}

func TestAtob(t *testing.T) {
	if !bytes.Equal(byteconv.Atob("test"), []byte("test")) {
		t.Error("should be equal")
	}
}
