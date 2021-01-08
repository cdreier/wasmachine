package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cdreier/wasmachine"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var images []image.Image

func main() {

	images = make([]image.Image, 0)

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Get("/api/generate", func(w http.ResponseWriter, r *http.Request) {

		// TODO

		if len(images) > 0 {
			w.Header().Add("Content-Type", "image/png")
			png.Encode(w, images[0])
			return
		}

	})

	r.Post("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("err in read body ", err)
		}

		img, _, err := image.Decode(bytes.NewReader(b))

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		images = append(images, img)

		json.NewEncoder(w).Encode(images)
	})

	if err := wasmachine.ListenAndServe(":1234", r); err != nil {
		fmt.Println("cannot start wasmachine", err.Error())
	}

}
