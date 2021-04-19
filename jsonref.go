package jsonref

import (
	"io"
)

type JsonRef interface {
	WriteTo(w io.Writer) (int64, error)
}
