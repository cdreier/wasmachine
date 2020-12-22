class Wasmachine {

  registration = null

  constructor(registration) {
    this.registration = registration
    this.registration.active.postMessage({
      type: "init"
    });
    navigator.serviceWorker.addEventListener('message', event => this.onMessage(event));
  }

  onMessage(evt) {
    switch (event.data.type) {
      case "fetch":
        this.request(JSON.stringify(event.data.req))
        break
      default:
        console.log(`JS: unhandled message received: ${event.data}`);
    }
  }

  request(req) {
    document.dispatchEvent(new CustomEvent("request", {
      detail: {
        req,
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