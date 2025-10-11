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

		for range 100 {
			v := pool.Get()
			v.String = "test"
			v.Pointer = &v.String
			pool.Put(v)
		}

		for range 100 {
			if !reflect.DeepEqual(&zero, pool.Get()) {
				t.Fatal("should be zero value")
			}
		}
	})

	t.Run("Pointer", func(t *testing.T) {
		zero := &item{}
		pool := don.NewRequestPool(zero)

		for range 100 {
			p := pool.Get()
			v := *p
			v.String = "test"
			v.Pointer = &v.String
			pool.Put(p)
		}

		for range 100 {
			if !reflect.DeepEqual(&zero, pool.Get()) {
				t.Fatal("should be zero value")
			}
		}
	})
}
