package extensions

import (
	"io"
	"net/http"

	"github.com/dihedron/rawdata"
)

type Response struct {
	URL     string              `json:"url"`
	Code    int                 `json:"code"`
	Status  string              `json:"status"`
	Headers map[string][]string `json:"headers"`
	Payload any                 `json:"payload"`
}

func CallAPI(url string) (*Response, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	payload, err := rawdata.Unmarshal(string(body))
	if err != nil {
		return nil, err
	}

	return &Response{
		URL:     url,
		Code:    response.StatusCode,
		Status:  response.Status,
		Headers: response.Header,
		Payload: payload,
	}, nil
}
