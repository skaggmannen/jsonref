package jsonref

import (
	"bytes"
	"io"
	"reflect"
)

type listRef struct {
	opts Opts
	t    reflect.Type
	v    reflect.Value
}

func (r listRef) WriteTo(w io.Writer) (int64, error) {
	var v reflect.Value
	if r.v.Len() > 0 {
		v = r.v.Field(0)
	}

	opts := r.opts

	var buf bytes.Buffer
	buf.WriteString("[ ")
	_, _ = kindRef{Opts: opts, t: r.t.Elem(), v: v}.WriteTo(&buf)
	buf.WriteString(" ... ]")

	return io.Copy(w, &buf)
}
