package utils

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
}

type HttpResponse struct {
	Code   int
	Body   []byte
	Header http.Header
}

func MakeRequest(data HttpRequest) (HttpResponse, error) {
	if data.Timeout == 0 {
		data.Timeout = 60 * time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), data.Timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, data.Method, data.Url, bytes.NewBuffer(data.Body))
	if err != nil {
		return HttpResponse{}, err
	}

	for k, v := range data.Header {
		request.Header.Set(k, v[0])
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return HttpResponse{}, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return HttpResponse{}, err
	}

	return HttpResponse{response.StatusCode, content, response.Header}, nil
}
