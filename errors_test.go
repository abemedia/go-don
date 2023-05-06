package don_test

import (
	"encoding/xml"
	"errors"
	"fmt"
	"testing"

	"github.com/abemedia/go-don"
	_ "github.com/abemedia/go-don/encoding/json"
	_ "github.com/abemedia/go-don/encoding/text"
	_ "github.com/abemedia/go-don/encoding/xml"
	_ "github.com/abemedia/go-don/encoding/yaml"
	"github.com/goccy/go-json"
	"github.com/google/go-cmp/cmp"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v3"
)

func TestError_Is(t *testing.T) {
	if !errors.Is(don.Error(errTest, 0), errTest) {
		t.Error("should match wrapped error")
	}
	if !errors.Is(don.Error(errTest, fasthttp.StatusBadRequest), don.ErrBadRequest) {
		t.Error("should match status error")
	}
}

func TestError_Unwrap(t *testing.T) {
	if errors.Unwrap(don.Error(errTest, 0)) != errTest {
		t.Error("should unwrap wrapped error")
	}
}

func TestError_StatusCode(t *testing.T) {
	if don.Error(don.ErrBadRequest, 0).StatusCode() != fasthttp.StatusBadRequest {
		t.Error("should respect wrapped error's status code")
	}
	if don.Error(don.ErrBadRequest, fasthttp.StatusConflict).StatusCode() != fasthttp.StatusConflict {
		t.Error("should ignore wrapped error's status code if explicitly set")
	}
}

func TestError_MarshalText(t *testing.T) {
	b, _ := don.Error(errTest, 0).MarshalText()
	if diff := cmp.Diff("test", string(b)); diff != "" {
		t.Error(diff)
	}
	b, _ = don.Error(&testError{}, 0).MarshalText()
	if diff := cmp.Diff("custom", string(b)); diff != "" {
		t.Error(diff)
	}
}

func TestError_MarshalJSON(t *testing.T) {
	b, _ := json.Marshal(don.Error(errTest, 0))
	if diff := cmp.Diff(`{"message":"test"}`, string(b)); diff != "" {
		t.Error(diff)
	}
	b, _ = json.Marshal(don.Error(&testError{}, 0))
	if diff := cmp.Diff(`{"custom":"test"}`, string(b)); diff != "" {
		t.Error(diff)
	}
}

func TestError_MarshalXML(t *testing.T) {
	b, _ := xml.Marshal(don.Error(errTest, 0))
	if diff := cmp.Diff("<message>test</message>", string(b)); diff != "" {
		t.Error(diff)
	}
	b, _ = xml.Marshal(don.Error(&testError{}, 0))
	if diff := cmp.Diff("<custom>test</custom>", string(b)); diff != "" {
		t.Error(diff)
	}
}

func TestError_MarshalYAML(t *testing.T) {
	b, _ := yaml.Marshal(don.Error(errTest, 0))
	if diff := cmp.Diff("message: test\n", string(b)); diff != "" {
		t.Error(diff)
	}
	b, _ = yaml.Marshal(don.Error(&testError{}, 0))
	if diff := cmp.Diff("custom: test\n", string(b)); diff != "" {
		t.Error(diff)
	}
}

var errTest = errors.New("test")

type testError struct{}

func (e *testError) Error() string {
	return "test"
}

func (e *testError) MarshalText() ([]byte, error) {
	return []byte("custom"), nil
}

func (e *testError) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"custom":%q}`, e.Error())), nil
}

func (e *testError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "custom"}
	return enc.EncodeElement(e.Error(), start)
}

func (e *testError) MarshalYAML() (any, error) {
	return map[string]string{"custom": e.Error()}, nil
}
