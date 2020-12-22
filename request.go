package wasmachine

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type FetchRequest struct {
	Cache               string            `json:"cache"`
	Credentials         string            `json:"credentials"`
	Destination         string            `json:"destination"`
	Headers             map[string]string `json:"headers"`
	Integrity           string            `json:"integrity"`
	IsHistoryNavigation bool              `json:"isHistoryNavigation"`
	Keepalive           bool              `json:"keepalive"`
	Method              string            `json:"method"`
	Mode                string            `json:"mode"`
	Redirect            string            `json:"redirect"`
	Referrer            string            `json:"referrer"`
	ReferrerPolicy      string            `json:"referrerPolicy"`
	URL                 string            `json:"url"`
	Body                string            `json:"body"`
	FetchID             string            `json:"fetchID"`
}

func NewRequest(fr FetchRequest) (*http.Request, error) {
	// path
	u, err := url.Parse(fr.URL)
	if err != nil {
		return nil, errors.Wrap(err, "parsing path")
	}

	// querystring
	q := u.Query()
	u.RawQuery = q.Encode()

	body := fr.Body

	// new request
	req, err := http.NewRequest(fr.Method, u.String(), strings.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	// manually set RequestURI because NewRequest is for clients and req.RequestURI is for servers
	req.RequestURI = u.RequestURI()

	// remote addr
	// TODO
	// req.RemoteAddr = e.RequestContext.Identity.SourceIP

	// header fields
	for k, v := range fr.Headers {
		req.Header.Set(k, v)
	}

	// content-length
	if req.Header.Get("Content-Length") == "" && body != "" {
		req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	}

	// host
	req.URL.Host = req.Header.Get("Host")
	req.Host = req.URL.Host

	return req, nil
}
