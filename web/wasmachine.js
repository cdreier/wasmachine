class Wasmachine {

  registration = null

  constructor(registration) {
    this.registration = registration
    navigator.serviceWorker.addEventListener('message', event => this.onMessage(event));
  }

  onMessage(evt) {
    switch (event.data.type) {
      case "fetch":
        this.request(event.data.req)
        break
      default:
        console.log(`JS: unhandled message received: ${event.data}`);
    }
  }

  request(req) {
    const requestDetail = Object.assign({}, req, {
      body: Array.from(new Uint8Array(req.body)),
    })

    document.dispatchEvent(new CustomEvent("request", {
      detail: {
        req: JSON.stringify(requestDetail),
      }
    }))
  }

  response(data) {
    this.registration.active.postMessage({
      type: "response",
      data,
    });
  }

}

window.onload = () => {
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('./worker.js').then(function (registration) {
      console.log('JS: ServiceWorker registration', registration.scope);
      document.workerBridge = new Wasmachine(registration)
    }).catch(function (err) {
      console.log('ServiceWorker registration failed: ', err);
    });
  }
}