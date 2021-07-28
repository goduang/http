package http

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpRequest struct {
	Method  string
	Url     string
	Body    []byte
	Header  http.Header
	Timeout time.Duration
	Client  *http.Client
}

type HttpResponse struct {
	Code   int
	Body   []byte
	Header http.Header
}

func MakeRequest(ctx context.Context, data *HttpRequest) (*HttpResponse, error) {
	request, err := http.NewRequestWithContext(ctx, data.Method, data.Url, bytes.NewBuffer(data.Body))
	if err != nil {
		return nil, err
	}

	for key, values := range data.Header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	if data.Client == nil {
		data.Client = http.DefaultClient
	}
	if data.Client.Timeout == 0 {
		data.Client.Timeout = data.Timeout
	}

	response, err := data.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{response.StatusCode, content, response.Header}, nil
}
