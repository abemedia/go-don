package protobuf

import (
	"reflect"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/encoding"
	"google.golang.org/protobuf/proto"
)

var messageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

func unmarshal(b []byte, v any) error {
	// TODO: Cache reflect results to improve performance.
	elem := reflect.ValueOf(v).Elem()
	if elem.Type().Implements(messageType) {
		v = elem.Interface()
	}

	m, ok := v.(proto.Message)
	if !ok {
		return don.ErrUnsupportedMediaType
	}
	return proto.Unmarshal(b, m)
}

func marshal(v any) ([]byte, error) {
	m, ok := v.(proto.Message)
	if !ok {
		return nil, don.ErrNotAcceptable
	}
	return proto.Marshal(m)
}

func init() {
	mediaType := "application/protobuf"
	alias := "application/x-protobuf"

	encoding.RegisterDecoder(unmarshal, mediaType, alias)
	encoding.RegisterEncoder(marshal, mediaType, alias)
}
