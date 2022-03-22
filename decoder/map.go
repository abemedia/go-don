package decoder

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
