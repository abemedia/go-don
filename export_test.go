package don

var NewNilCheck = newNilCheck

func NewRequestPool[T any](v T) pool[T] {
	return newRequestPool(v)
}
