package decoder

import (
	"github.com/abemedia/go-don/internal/byteconv"
	"github.com/abemedia/httprouter"
	"github.com/valyala/fasthttp"
)

type Map map[string][]string

func (m Map) Get(key string) string {
	if m == nil {
		return ""
	}
	if vs := m[key]; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (m Map) Values(key string) []string {
	if m == nil {
		return nil
	}
	return m[key]
}

type Args fasthttp.Args

func (ps *Args) Get(key string) string {
	return byteconv.Btoa((*fasthttp.Args)(ps).Peek(key))
}

func (ps *Args) Values(key string) []string {
	args := (*fasthttp.Args)(ps).PeekMulti(key)
	if len(args) == 0 {
		return nil
	}

	res := make([]string, len(args))
	for i := range args {
		res[i] = byteconv.Btoa(args[i])
	}

	return res
}

type Header fasthttp.RequestHeader

func (ps *Header) Get(key string) string {
	return byteconv.Btoa((*fasthttp.RequestHeader)(ps).Peek(key))
}

func (ps *Header) Values(key string) []string {
	args := (*fasthttp.RequestHeader)(ps).PeekAll(key)
	if len(args) == 0 {
		return nil
	}

	res := make([]string, len(args))
	for i := range args {
		res[i] = byteconv.Btoa(args[i])
	}

	return res
}

type Params httprouter.Params

func (ps Params) Get(key string) string {
	for i := range ps {
		if ps[i].Key == key {
			return ps[i].Value
		}
	}
	return ""
}

func (ps Params) Values(key string) []string {
	for i := range ps {
		if ps[i].Key == key {
			return []string{ps[i].Value}
		}
	}
	return nil
}
