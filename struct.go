package jsonref

import (
	"bytes"
	"io"
	"reflect"
	"strings"
)

func Struct(s interface{}, structOpts ...StructOpt) JsonRef {
	t := reflect.TypeOf(s)

	if t.Kind() != reflect.Struct {
		panic("input must be a struct")
	}

	opts := Opts{}
	for _, o := range structOpts {
		o.Apply(&opts)
	}

	return structRef{Opts: opts, t: t, v: reflect.ValueOf(s)}
}

type StructOpt interface {
	Apply(opts *Opts)
}

type StructOptFunc func(o *Opts)

func (o StructOptFunc) Apply(opts *Opts) {
	o(opts)
}

func HrefSep(s string) StructOptFunc {
	return func(o *Opts) {
		o.HrefSep = s
	}
}

func Ignore(paths ...string) StructOptFunc {
	return func(o *Opts) {
		o.Ignore = append(o.Ignore, paths...)
	}
}

type Lookup map[string][]string

func (l Lookup) Apply(o *Opts) {
	o.Lookup = l
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

	if opts.ignored() {
		return 0, nil
	}

	var buf bytes.Buffer
	buf.WriteString(opts.indent() + opts.fieldName(name))

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
		if values, found := r.Lookup[oneOf]; found {
			return values
		}

		return strings.Split(oneOf, ",")
	}

	return nil
}
