self.addEventListener('message', msg => {

  switch (msg.data.type) {
    case "init":
      self._parent = msg.source
      self.requestBuffer = {}
      break;
    case "response":
      // console.log("SW: response", msg.data)
      const res = JSON.parse(msg.data.data)
      const { done } = self.requestBuffer[res.fetchID]

      // TODO build more attributes in response ?
      const headers = new Headers()
      Object.keys(res.headers).forEach(k => {
        headers.append(k, res.headers[k])
      })

      const responseBody = res.isBinary ? parseBinaryResponse(res.body) : res.body

      done(new Response(responseBody, {
        status: res.statusCode,
        headers,
      }))

      delete self.requestBuffer[res.fetchID]
      break;
  }
});

const parseBinaryResponse = (res) => {
  const byteChars = atob(res)
  const byteNumbers = new Array(byteChars.length);
  for (let i = 0; i < byteChars.length; i++) {
    byteNumbers[i] = byteChars.charCodeAt(i);
  }
  const byteArray = new Uint8Array(byteNumbers);
  const blob = new Blob([byteArray]);
  return blob
}

const serializeRequest = async (req, fetchID) => {
  // TODO: perhaps we need to check the body 
  // https://developer.mozilla.org/en-US/docs/Web/API/Request#Methods
  const body = await req.arrayBuffer()
  const r = {
    cache: req.cache,
    credentials: req.credentials,
    destination: req.destination,
    headers: Object.fromEntries(req.headers),
    integrity: req.integrity,
    isHistoryNavigation: req.isHistoryNavigation,
    keepalive: req.keepalive,
    method: req.method,
    mode: req.mode,
    redirect: req.redirect,
    referrer: req.referrer,
    referrerPolicy: req.referrerPolicy,
    url: req.url,
    fetchID,
    body,
  }
  return r
}

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
      }).catch(err => console.log(err))
    )
  }
});

function uuidv4() {
  return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
  );
}
