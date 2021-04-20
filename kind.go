package jsonref

import (
	"io"
	"reflect"
	"strings"
)

type kindRef struct {
	Opts
	v reflect.Value
	t reflect.Type
}

func (r kindRef) WriteTo(w io.Writer) (int64, error) {
	switch r.t.Kind() {

	case reflect.String:
		return stringRef{
			Opts: r.Opts,
			t:    r.t,
			v:    r.v,
		}.WriteTo(w)
	case reflect.Interface:
		return interfaceRef{
			Opts: r.Opts,
			v:    r.v,
			t:    r.t,
		}.WriteTo(w)
	case reflect.Array, reflect.Slice:
		return listRef{
			Opts: r.Opts,
			t:    r.t,
			v:    r.v,
		}.WriteTo(w)
	case reflect.Struct:
		return structRef{
			Opts: r.Opts,
			t:    r.t,
			v:    r.v,
		}.WriteTo(w)
	case reflect.Ptr:
		return ptrRef{
			Opts: r.Opts,
			t:    r.t,
			v:    r.v,
		}.WriteTo(w)
	default:
		t := r.v.Kind().String()

		if r.Format != "" {
			t = r.Format
		}

		return io.Copy(w, strings.NewReader(r.fieldType(t)))
	}
}
