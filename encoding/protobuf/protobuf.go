// Package protobuf provides encoding and decoding of Protocol Buffers data.
package protobuf

import (
	"reflect"
	"sync"

	"github.com/abemedia/go-don"
	"github.com/abemedia/go-don/encoding"
	"google.golang.org/protobuf/proto"
)

var (
	messageType = reflect.TypeFor[proto.Message]()
	cache       sync.Map
)

func unmarshal(b []byte, v any) error {
	typ := reflect.TypeOf(v)
	fn, ok := cache.Load(typ)
	if !ok {
		if typ.Elem().Implements(messageType) {
			fn = func(v any) any { return reflect.ValueOf(v).Elem().Interface() }
		} else {
			fn = func(v any) any { return v }
		}
		cache.Store(typ, fn)
	}

	m, ok := fn.(func(v any) any)(v).(proto.Message)
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
