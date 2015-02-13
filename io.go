package feedback

import (
	"io"
)

type readWriter struct {
	io.Reader
	io.Writer
}

// NewReadWriter combines a reader and writer into an io.ReadWriter
// This really should be part of the io package...
func NewReadWriter(r io.Reader, w io.Writer) io.ReadWriter {
	return &readWriter{r, w}
}
