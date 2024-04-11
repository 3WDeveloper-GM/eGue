package adapter

import (
	"io"
	"net/http"

	"golang.org/x/xerrors"
)

var contentType = "application/json"
var usrAgnt = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"

type SearchAdapter struct {
	cfg    DBImplementation
	client *http.Client
}

func NewAdapter(client *http.Client, config DBImplementation) *SearchAdapter {

	return &SearchAdapter{
		client: client,
		cfg:    config,
	}
}

func (za *SearchAdapter) Generate(method string, url string, body io.Reader) (*http.Request, error) {
	var request *http.Request
	var err error
	switch {
	case method == http.MethodPost:

		request, err = http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}
		// user-agents and content type for interacting with zincSearch
		request.Header.Set("Content-Type", contentType)
		request.Header.Set("User-Agent", usrAgnt)

		request.SetBasicAuth(za.cfg.GetDBAdmin(), za.cfg.GetDBPassword())

	case method == http.MethodGet:

		request, err = http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}

		request.Header.Set("Content-Type", contentType)
		request.Header.Set("User-Agent", usrAgnt)

		request.SetBasicAuth(za.cfg.GetDBAdmin(), za.cfg.GetDBPassword())

	default:
		err = xerrors.Errorf("method not valid: only POST methods are used")
		return nil, err
	}
	return request, nil
}

func (za *SearchAdapter) Do(request *http.Request) (*http.Response, error) {
	return za.client.Do(request)
}
