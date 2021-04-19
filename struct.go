package jsonref

import (
	"bytes"
	"io"
	"reflect"
	"strings"
)

func Struct(s interface{}, optFuns ...OptFunc) JsonRef {
	t := reflect.TypeOf(s)

	if t.Kind() != reflect.Struct {
		panic("input must be a struct")
	}

	opts := Opts{}
	for _, o := range optFuns {
		o(&opts)
	}

	return structRef{Opts: opts, t: t, v: reflect.ValueOf(s)}
}

type OptFunc func(o *Opts)

func HrefSep(s string) OptFunc {
	return func(o *Opts) {
		o.HrefSep = s
	}
}

type structRef struct {
	Opts
	t reflect.Type
	v reflect.Value
}

func (r structRef) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer
	buf.WriteString("{\n")

	for i := 0; i < r.t.NumField(); i++ {
		var v reflect.Value
		if r.v.IsValid() {
			v = r.v.Field(i)
		}

		_, _ = structFieldRef{
			Opts: r.Opts,
			t:    r.t.Field(i),
			v:    v,
		}.WriteTo(&buf)
	}

	buf.WriteString(r.indent() + "}")

	return io.Copy(w, &buf)
}

type structFieldRef struct {
	Opts
	t reflect.StructField
	v reflect.Value
}

func (r structFieldRef) WriteTo(w io.Writer) (int64, error) {
	name := r.name()
	if name == "-" {
		return 0, nil
	}

	opts := r.Opts
	opts.Explicit = r.explicit()
	opts.Format = r.format()
	opts.OneOf = r.oneOf()
	opts.Parents = append(opts.Parents, name)
	opts.IndentLevel += 1

	var buf bytes.Buffer
	buf.WriteString(opts.indent() + r.fieldName(name))

	_, _ = kindRef{
		Opts: opts,
		v:    r.v,
		t:    r.t.Type,
	}.WriteTo(&buf)

	buf.WriteString(",\n")

	return io.Copy(w, &buf)
}

func (r structFieldRef) explicit() bool {
	_, ok := r.t.Tag.Lookup("explicit")
	return ok
}

func (r structFieldRef) name() string {
	name := r.t.Name

	if jsonTag, ok := r.t.Tag.Lookup("json"); ok {
		i := strings.Index(jsonTag, ",")
		if i < 0 {
			name = jsonTag
		} else {
			name = jsonTag[:i]
		}
	}

	return name
}

func (r structFieldRef) format() string {
	if format, ok := r.t.Tag.Lookup("format"); ok {
		return format
	}

	return ""
}

func (r structFieldRef) oneOf() []string {
	if oneOf, ok := r.t.Tag.Lookup("oneOf"); ok {
		return strings.Split(oneOf, ",")
	}

	return nil
}