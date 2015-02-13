package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	feedback "github.com/jbenet/go-feedback"
)

var (
	listenAddr string
)

func init() {
	flag.StringVar(&listenAddr, "listen", ":0", "listen address")
	flag.StringVar(&listenAddr, "l", ":0", "listen address (shorthand)")
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	l, err := net.Listen("tcp4", listenAddr)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Listening at %s\n", l.Addr())
	Serve(l, os.Stdout)
	return nil
}

func Serve(l net.Listener, output io.Writer) {

	handler := func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "POST" {
			return fmt.Errorf("invalid method")
		}

		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}

		f, err := feedback.Unmarshal(buf)
		if err != nil {
			return err
		}

		buf2, err := feedback.Marshal(f)
		if err != nil {
			return err
		}

		output.Write(buf2)
		output.Write([]byte("\n"))
		w.WriteHeader(http.StatusOK)
		w.Write(buf2)
		return nil
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	})

	http.Serve(l, nil)
}
