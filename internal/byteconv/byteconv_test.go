package byteconv_test

import (
	"bytes"
	"testing"

	"github.com/abemedia/go-don/internal/byteconv"
)

func TestBtoa(t *testing.T) {
	b := []byte("test")
	if string(b) != byteconv.Btoa(b) {
		t.Error("should be equal")
	}
}

func TestAtob(t *testing.T) {
	s := "test"
	if !bytes.Equal([]byte(s), byteconv.Atob(s)) {
		t.Error("should be equal")
	}
}
