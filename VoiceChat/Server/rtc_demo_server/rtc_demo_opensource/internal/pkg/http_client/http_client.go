package http_client

import (
	"github.com/valyala/fasthttp"
)

type HeaderPair struct {
	Key   string
	Value string
}

func NewHeaderPair(key, value string) HeaderPair {
	return HeaderPair{
		Key:   key,
		Value: value,
	}
}

var client = &fasthttp.Client{}

func DoRequest(url, method string, body []byte, headers ...HeaderPair) (int, []byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetRequestURI(url)
	req.SetBody(body)
	req.Header.SetMethod(method)
	req.Header.Set("Content-Type", "application/json")

	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}

	if err := client.Do(req, resp); err != nil {
		return 0, nil, err
	}

	return resp.StatusCode(), resp.Body(), nil
}
