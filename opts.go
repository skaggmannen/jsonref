package jsonref

import (
	"fmt"
	"strings"
)

type Opts struct {
	Explicit    bool
	Format      string
	OneOf       []string
	IndentLevel int
	Parents     []string
	HrefSep     string
	Ignore      []string

	Fmt FmtOpts
}

type FmtOpts struct {
	Href       func(parents []string) string
	Indent     func(level int) string
	FieldName  func(name string, href string) string
	FieldType  func(t string) string
	FieldValue func(v string) string
}

func (o Opts) indent() string {
	if o.Fmt.Indent != nil {
		return o.Fmt.Indent(o.IndentLevel)
	}

	return strings.Repeat("  ", o.IndentLevel)
}

func (o Opts) href() string {
	if o.Fmt.Href != nil {
		return o.Fmt.Href(o.Parents)
	}

	return "#" + strings.ToLower(strings.Join(o.Parents, o.HrefSep))
}

func (o Opts) fieldName(name string) string {
	if o.Fmt.FieldName != nil {
		return o.Fmt.FieldName(name, o.href())
	}

	return fmt.Sprintf(`<a href="%s" class="jsonref-link">%s</a>: `, o.href(), name)
}

func (o Opts) fieldType(t string) string {
	if o.Fmt.FieldType != nil {
		return o.Fmt.FieldType(t)
	}

	return `<span class="jsonref-type">` + t + `</span>`
}

func (o Opts) fieldValue(v string) string {
	if o.Fmt.FieldValue != nil {
		return o.Fmt.FieldValue(v)
	}

	return `<span class="jsonref-value">` + v + `</span>`
}

func (o Opts) ignored() bool {
	path := strings.ToLower(strings.Join(o.Parents, "."))

	for _, ignored := range o.Ignore {
		if path == strings.ToLower(ignored) {
			return true
		}
	}

	return false
}
