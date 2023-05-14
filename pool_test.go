package don_test

import (
	"reflect"
	"testing"

	"github.com/abemedia/go-don"
)

func TestRequestPool(t *testing.T) {
	type item struct {
		String  string
		Pointer *string
	}

	t.Run("Nil", func(t *testing.T) {
		var zero any
		pool := don.NewRequestPool(zero)

		pool.Put(pool.Get())

		if !reflect.DeepEqual(&zero, pool.Get()) {
			t.Fatal("should be zero value")
		}
	})

	t.Run("Struct", func(t *testing.T) {
		zero := item{}
		pool := don.NewRequestPool(zero)

		for i := 0; i < 100; i++ {
			v := pool.Get()
			v.String = "test"
			v.Pointer = &v.String
			pool.Put(v)
		}

		for i := 0; i < 100; i++ {
			if !reflect.DeepEqual(&zero, pool.Get()) {
				t.Fatal("should be zero value")
			}
		}
	})

	t.Run("Pointer", func(t *testing.T) {
		zero := &item{}
		pool := don.NewRequestPool(zero)

		for i := 0; i < 100; i++ {
			p := pool.Get()
			v := *p
			v.String = "test"
			v.Pointer = &v.String
			pool.Put(p)
		}

		for i := 0; i < 100; i++ {
			if !reflect.DeepEqual(&zero, pool.Get()) {
				t.Fatal("should be zero value")
			}
		}
	})

	t.Run("Resetter", func(t *testing.T) {
		zero := &itemResetter{}
		pool := don.NewRequestPool(zero)

		for i := 0; i < 100; i++ {
			p := pool.Get()
			v := *p
			v.String = "test"
			v.Pointer = &v.String
			pool.Put(p)
		}

		for i := 0; i < 100; i++ {
			if !reflect.DeepEqual(&zero, pool.Get()) {
				t.Fatal("should be zero value")
			}
		}
	})
}

type itemResetter struct {
	String  string
	Pointer *string
}

func (ir *itemResetter) Reset() {
	*ir = itemResetter{}
}
