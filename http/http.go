package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient interface {
	Get(url string) ([]byte, error)
}

type JsonUnmarshaller interface {
	Unmarshal(json []byte, v interface{}) error
}

type httpClient struct {
	timeout time.Duration
}

func (c *httpClient) Get(url string) ([]byte, error) {
	client := http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func NewHttpClient(timeout time.Duration) HttpClient {
	return &httpClient{timeout: timeout}
}

type jsonUnmarshaller struct {
}

func (ju *jsonUnmarshaller) Unmarshal(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}

func NewJsonUnmarshaller() JsonUnmarshaller {
	return &jsonUnmarshaller{}
}
