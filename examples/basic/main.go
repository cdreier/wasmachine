package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cdreier/wasmachine"
	"github.com/go-chi/chi"
)

func main() {

	r := chi.NewRouter()

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TEST GET from wasm"))
	})

	r.Post("/api/debug", func(w http.ResponseWriter, r *http.Request) {

		b := struct {
			ServerTime    time.Time `json:"server_time,omitempty"`
			ServerMessage string    `json:"server_message,omitempty"`
			Txt           string    `json:"txt,omitempty"`
		}{}
		json.NewDecoder(r.Body).Decode(&b)

		b.ServerTime = time.Now()
		b.ServerMessage = "nice! now i only need a usecase :)"

		json.NewEncoder(w).Encode(b)
	})

	if err := wasmachine.ListenAndServe(":1234", r); err != nil {
		fmt.Println("cannot start wasmachine", err.Error())
	}

}
