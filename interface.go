package jsonref

import (
	"io"
	"reflect"
	"strings"
)

type interfaceRef struct {
	Opts
	t reflect.Type
	v reflect.Value
}

func (r interfaceRef) WriteTo(w io.Writer) (int64, error) {
	if r.v.IsNil() {
		return io.Copy(w, strings.NewReader(r.fieldType("Unknown")))
	}

	v := r.v.Elem()

	return kindRef{Opts: r.Opts, t: v.Type(), v: v}.WriteTo(w)
}
