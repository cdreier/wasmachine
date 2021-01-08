package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cdreier/wasmachine"
	"github.com/go-chi/chi"
)

func main() {

	r := chi.NewRouter()

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TEST GET from wasm"))
	})

	r.Post("/api/debug", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("err in read body ", err)
		}
		w.Write([]byte(fmt.Sprintf("TEST POST from wasm: %s", string(b))))
	})

	if err := wasmachine.ListenAndServe(":1234", r); err != nil {
		fmt.Println("cannot start wasmachine", err.Error())
	}

}
