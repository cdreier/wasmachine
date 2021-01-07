package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cdreier/wasmachine"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TEST GET from wasm"))
	})

	r.Post("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("err in read body ", err)
		}

		img, format, err := image.Decode(bytes.NewReader(b))

		if err != nil {
			log.Println(err, format)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println(img.Bounds(), format)

		w.Write([]byte("upload endpoint"))
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
