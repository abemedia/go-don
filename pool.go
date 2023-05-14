package don

import (
	"reflect"
	"sync"
)

type pool[T any] interface {
	Get() *T
	Put(*T)
}

type resetter interface {
	Reset()
}

var resetterType = reflect.TypeOf((*resetter)(nil)).Elem()

type requestPool[T any] struct {
	pool  sync.Pool
	reset func(*T)
}

func newRequestPool[T any](zero T) pool[T] {
	typ := reflect.TypeOf(zero)
	if typ == nil {
		return &fakePool[T]{&zero}
	}

	p := &requestPool[T]{}

	if typ.Kind() != reflect.Pointer {
		p.pool.New = func() any {
			return new(T)
		}
		p.reset = func(v *T) {
			*v = zero
		}
	} else {
		elem := typ.Elem()
		p.pool.New = func() any {
			v := reflect.New(elem).Interface().(T) //nolint:forcetypeassert
			return &v
		}

		if typ.Implements(resetterType) {
			p.reset = func(v *T) {
				any(*v).(resetter).Reset() //nolint:forcetypeassert
			}
		} else {
			zeroValue := reflect.New(elem).Elem()
			p.reset = func(v *T) {
				reflect.ValueOf(v).Elem().Elem().Set(zeroValue)
			}
		}
	}

	return p
}

func (p *requestPool[T]) Get() *T {
	return p.pool.Get().(*T) //nolint:forcetypeassert
}

func (p *requestPool[T]) Put(v *T) {
	p.reset(v)
	p.pool.Put(v)
}

type fakePool[T any] struct{ v *T }

func (p *fakePool[T]) Get() *T { return p.v }

func (p *fakePool[T]) Put(*T) {}
