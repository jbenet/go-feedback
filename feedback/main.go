package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	feedback "github.com/jbenet/go-feedback"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	rw := NewReadWriter(os.Stdin, os.Stdout)

	f, err := feedback.PromptForFeedback(rw, feedback.Options{})
	if err != nil {
		return err
	}

	buf, err := json.Marshal(&f)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(buf)
	return err
}

type readWriter struct {
	io.Reader
	io.Writer
}

func NewReadWriter(r io.Reader, w io.Writer) io.ReadWriter {
	return &readWriter{r, w}
}
