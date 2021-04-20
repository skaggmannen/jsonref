package jsonref_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/skaggmannen/jsonref"
)

func Test_structRef(t *testing.T) {
	t.Run("example struct", func(t *testing.T) {
		var buf bytes.Buffer

		ref := jsonref.Struct(
			TestStruct{
				A: "example-value",
				D: TestStruct{},
			},
			jsonref.HrefSep("-"),
			jsonref.Ignore("d.m.a"),
			jsonref.Lookup{
				"allowedValues": []string{"Foo", "Bar", "Baz"},
			},
		)

		_, _ = ref.WriteTo(&buf)

		s := buf.String()
		log.Printf(s)

		if s == "" {
			t.Error("Output should not empty")
		}
	})
}

type TestStruct struct {
	A string      `json:"a"`
	B string      `json:"b" oneOf:"allowedValues"`
	C string      `json:"-"`
	D interface{} `json:"d" format:"shouldNotBeDisplayed"`
	E string      `json:"e" format:"UUID"`
	F string      `json:"f" oneOf:"A,B,C,D"`
	G [][]byte    `json:"g" format:"Base64"`
	H int8        `json:"h"`
	I uint8       `json:"i" format:"Hour"`
	J uint8       `json:"j" format:"Minute"`
	K uint8       `json:"k" format:"Second"`
	L float64     `json:"l"`
	M []struct {
		A bool `json:"a"`
	} `json:"m"`
	N []byte `json:"n"`
}
