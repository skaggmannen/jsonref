package jsonref

import (
	"io"
	"reflect"
	"strings"
)

type stringRef struct {
	Opts
	t reflect.Type
	v reflect.Value
}

func (r stringRef) WriteTo(w io.Writer) (int64, error) {
	t := r.fieldType("String")

	if r.Format != "" {
		t = r.fieldType(r.Format)
	}

	if len(r.OneOf) > 0 {
		t = r.fieldValue(`"` + strings.Join(r.OneOf, " | ") + `"`)
	}

	if r.v.IsValid() {
		s := r.v.String()
		if len(s) > 0 {
			t = r.fieldValue(`"` + s + `"`)
		}

	}

	return io.Copy(w, strings.NewReader(t))
}
