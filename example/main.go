package main

import (
	"fmt"
	"net/http"

	"github.com/cdreier/wasmachine"
	"github.com/go-chi/chi"
)

func main() {

	r := chi.NewRouter()

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GET from wasm"))
	})

	r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("POST from wasm"))
	})

	if err := wasmachine.ListenAndServe(":1234", r); err != nil {
		fmt.Println("cannot start wasmachine", err.Error())
	}

}
