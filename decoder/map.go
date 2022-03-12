package decoder

type MapDecoder struct {
	dec *CachedDecoder
}

func NewMapDecoder(v interface{}, tag string) (*MapDecoder, error) {
	dec, err := NewCachedDecoder(v, tag)
	if err != nil {
		return nil, err
	}

	return &MapDecoder{dec}, nil
}

func (d *MapDecoder) Decode(data map[string][]string, v interface{}) error {
	return d.dec.Decode(MapGetter(data), v)
}

type MapGetter map[string][]string

func (m MapGetter) Get(key string) string {
	if m == nil {
		return ""
	}

	vs := m[key]
	if len(vs) == 0 {
		return ""
	}

	return vs[0]
}

func (m MapGetter) Values(key string) []string {
	if m == nil {
		return nil
	}

	return m[key]
}
