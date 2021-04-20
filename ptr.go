package jsonref

import (
	"bytes"
	"io"
	"reflect"
)

type ptrRef struct {
	Opts
	t reflect.Type
	v reflect.Value
}

func (r ptrRef) WriteTo(w io.Writer) (int64, error) {
	t := r.t
	v := r.v

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		if v.IsValid() {
			v = v.Elem()
		}
	}

	var buf bytes.Buffer
	buf.WriteString("* ")
	_, _ = kindRef{
		Opts: r.Opts,
		t:    t,
		v:    v,
	}.WriteTo(&buf)

	return io.Copy(w, &buf)
}
