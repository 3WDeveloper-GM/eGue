package crawler

import (
	"io"
	"net/http"
)

type FileLogger interface {
	Log(message string)
}

type PostClient interface {
	Do(*http.Request) (*http.Response, error)
	Generate(method string, url string, reader io.Reader) (*http.Request, error)
}
