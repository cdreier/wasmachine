package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cdreier/wasmachine"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var frames []*image.Paletted

func main() {

	frames = make([]*image.Paletted, 0)

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Get("/api/generate", func(w http.ResponseWriter, r *http.Request) {

		delays := make([]int, len(frames))
		for i := range frames {
			delays[i] = 0
		}

		if len(frames) > 0 {
			w.Header().Add("Content-Type", "image/gif")
			gif.EncodeAll(w, &gif.GIF{
				Image: frames,
				Delay: delays,
			})
			return
		}

		w.Write([]byte("nothing to do..."))

	})

	r.Post("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("err in read body ", err)
		}

		img, err := gif.DecodeAll(bytes.NewReader(b))

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		frames = append(frames, img.Image...)

		json.NewEncoder(w).Encode(struct {
			Frames int `json:"frames,omitempty"`
		}{
			Frames: len(frames),
		})
	})

	if err := wasmachine.ListenAndServe(":1234", r); err != nil {
		fmt.Println("cannot start wasmachine", err.Error())
	}

}
