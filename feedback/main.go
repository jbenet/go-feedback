package main

import (
	"flag"
	"log"
	"os"

	feedback "github.com/jbenet/go-feedback"
)

var (
	flagURL string
)

func init() {
	flag.StringVar(&flagURL, "post", "", "POST results to given URL")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flag.Parse()

	rw := feedback.NewReadWriter(os.Stdin, os.Stderr)

	f, err := feedback.PromptForFeedback(rw, feedback.Options{})
	if err != nil {
		return err
	}

	if flagURL != "" {
		return feedback.PostFeedback(f, flagURL)
	}

	buf, err := feedback.Marshal(f)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(buf)
	return err
}
