# wasmachine

use your golang http router in webassembly.

It's as easy as 

1. code: `wasmachine.ListenAndServe(":1234", yourHTTPRouter)`
2. build: `GOOS=js GOARCH=wasm go build -o main.wasm`
3. ... adding a bunch of javascripts beside the wasm ;) 

## Examples

the targets in the Makefile are building the examples and copies everything needed in the `web` folder. Now just start a simple webserver with the `web` folder as root, and you are good to go!

You can find the gif demo running [here](https://drailing.net/demos/wasmachine/)!

## HOW???

This is a great question - glad you asked. 

Playing around with service workers, i found [FetchEvents](https://developer.mozilla.org/en-US/docs/Web/API/FetchEvent). The day before i played with go and the wasm output - wondering why there is no easy and convinient way to call a function from JS and just get the output. One thing came to another. 

Here we are now. 

With a wonderful **wasmachine**

### MORE DETAILS!

Alright - a common usecase for the service worker `FetchEvent` is to intercept and cache responses and resources. But service workers are in their own world, not access to the DOM or anything. 

To call something from the service working in the wasm, we need a bridge between the two worlds, gladly there is the [postMessage API](https://developer.mozilla.org/en-US/docs/Web/API/Client/postMessage)

#### MORE CODE!

let's start simple and just intercept the fetch request only on a specifiy path

```js
self.addEventListener('fetch', (event) => {
  if (event.request.url.includes("/api/")) {
    // yay!
  }
})
```

now the strange part... we need to serialize the request, `postMessage` it back to the DOM (as there is our wasm access). In the DOM context we custom- `dispatchEvent` the serialized request. With the golang `syscall/js` package, we register a callback to our custom event, named "request"

```go
js.Global().Get("document").Call("addEventListener", "request", cb)
```

Then i borrowed the idea from the [apex gateway](https://github.com/apex/gateway) - they create golang `http.Request`s from the AWS API Gateway so we can handle lambda requests like normal http requests (this is awesome btw!) - let's take a look at  `cb` - the callback function registered above.

```go
cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
    // get the payload of the custom event
    payload := args[0].Get("detail")
    // get the serialized request and unmarshal
    jsReq := payload.Get("req").String()
    var freq FetchRequest
    json.Unmarshal([]byte(jsReq), &freq)

    // create a go http.Request! -> request.go
    r, _ := NewRequest(freq) 
    // and a response -> response.go
    w := NewResponse(freq.FetchID)
    // and call some handler :) 
    h.ServeHTTP(w, r)

    resp := w.End()
    jsonResp, _ := json.Marshal(resp)

    // and call the response function in our javascript bridge
    js.Global().Get("document").Get("workerBridge").Call("response", string(jsonResp))
    return nil
})
```

wow.

And to use it like a real webserver, i put it in a `ListenAndServe` function, so you can just call

```go
r := chi.NewRouter()
// ... handlers
if err := wasmachine.ListenAndServe(":1234", r); err != nil {
  fmt.Println("cannot start wasmachine", err.Error())
}
```

The port is only there for the *feeling*.

When we put everything together, we see that we need to store the request somewhere in the service worker to answer it after our go wasmachine handled everything. This was a bit tricky. Luckily we can respond with a `Promise`! 

```js
self.addEventListener('fetch', (event) => {
  if (event.request.url.includes("/api/")) {
    event.respondWith(
      new Promise(async resolve => {
        const fetchID = uuidv4()
        const req = await serializeRequest(event.request, fetchID)
        self._parent.postMessage({
          type: "fetch",
          req: req,
        });
        self.requestBuffer[fetchID] = {
          done: resolve,
        }
      })
    )
  }
});
```

Our `requestBuffer` is just a map where we store a generated fetchID and the `resolve` function from the promise. After wasmachine posts back the response, we resolve the request, construct a JS response and we are done!

```js
const res = JSON.parse(msg.data.data)
const { done } = self.requestBuffer[res.fetchID]

const headers = new Headers()
Object.keys(res.headers).forEach(k => {
  headers.append(k, res.headers[k])
})
done(new Response(res.body, {
  status: res.statusCode,
  headers,
}))
```

The solution for binary payloads and responses is the same as AWS is doing it - pack everything in a base64 string and unpack it on the receiver.

## WHY???

I DONT KNOW - Y U YELLING AT ME

Perhaps we can concat gifs... or do something cool with a [SingleHostReverseProxy](https://golang.org/pkg/net/http/httputil/#NewSingleHostReverseProxy)?

This is just the beginning! ;) 