package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/romshark/templbench/templates"
)

func main() {
	fHost := flag.String("host", ":8080", "HTTP host address")
	flag.Parse()

	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/helloworld":
			if r.Method != http.MethodGet {
				const c = http.StatusMethodNotAllowed
				http.Error(w, http.StatusText(c), c)
				return
			}
			err := templates.RenderHelloWorld(
				r.Context(), w, "HelloWorld", "Hello World!",
			)
			if err != nil {
				panic(err)
			}
		default:
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	})

	if err := http.ListenAndServe(*fHost, handle); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}
}
