package jsonref

import (
	"bytes"
	"io"
	"reflect"
	"strings"
)

type listRef struct {
	Opts
	t reflect.Type
	v reflect.Value
}

func (r listRef) WriteTo(w io.Writer) (int64, error) {
	if r.t.Elem().Kind() == reflect.Uint8 {
		return io.Copy(w, strings.NewReader(r.fieldType("Base64")))
	}

	var v reflect.Value
	if v.IsValid() && r.v.Len() > 0 {
		v = r.v.Index(0)
	}

	opts := r.Opts

	var buf bytes.Buffer
	buf.WriteString("[ ")
	_, _ = kindRef{Opts: opts, t: r.t.Elem(), v: v}.WriteTo(&buf)
	buf.WriteString(" ... ]")

	return io.Copy(w, &buf)
}
