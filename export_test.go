package don

var MakeNilCheck = makeNilCheck

func NewRequestPool[T any](v T) pool[T] {
	return newRequestPool(v)
}
