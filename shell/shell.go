package shell

import (
	"io"
)

type Shell interface {
	Add(r io.Reader) (string, error)
	Cat(path string) (io.ReadCloser, error)
}
