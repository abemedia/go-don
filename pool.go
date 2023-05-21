package don

import (
	"reflect"
	"sync"
)

type pool[T any] interface {
	Get() *T
	Put(*T)
}

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
		rtype := dataOf(typ)
		elem := typ.Elem()
		elemrtype := dataOf(elem)
		zero := dataOf(reflect.New(elem).Elem().Interface())

		p.pool.New = func() any {
			v := packEface(rtype, unsafe_New(elemrtype)).(T) //nolint:forcetypeassert
			return &v
		}
		p.reset = func(v *T) {
			typedmemmove(elemrtype, dataOf(*v), zero)
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
