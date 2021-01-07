package wasmachine

import (
	"encoding/json"
	"log"
	"net/http"
	"syscall/js"
)

func ListenAndServe(addr string, h http.Handler) error {
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		payload := args[0].Get("detail")

		// test := make([]byte, 0)
		// fmt.Println("JS payload: ", payload.Call("toString"))
		// bytesCopied := js.CopyBytesToGo(test, payload)
		jsReq := payload.Get("req").String()
		var freq FetchRequest
		err := json.Unmarshal([]byte(jsReq), &freq)
		if err != nil {
			log.Println("error during req deserialization: ", err)
		}

		r, _ := NewRequest(freq)
		w := NewResponse(freq.FetchID)
		h.ServeHTTP(w, r)

		resp := w.End()
		jsonResp, _ := json.Marshal(resp)

		js.Global().Get("document").Get("workerBridge").Call("response", string(jsonResp))
		return nil
	})
	js.Global().Get("document").Call("addEventListener", "request", cb)

	done := make(chan bool, 1)

	<-done

	return nil
}
