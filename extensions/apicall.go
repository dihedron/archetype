package extensions

import (
	"io"
	"net/http"

	"github.com/dihedron/rawdata"
)

// Response is the structure returned by the CallAPI function.
// It contains the response URL, status code and message, headers and a
// payload, which is the result of unmarshalling the response body.
type Response struct {
	URL     string              `json:"url"`
	Code    int                 `json:"code"`
	Status  string              `json:"status"`
	Headers map[string][]string `json:"headers"`
	Payload any                 `json:"payload"`
}

// CallAPI is a template function that makes a GET request to the given URL
// and returns a Response object.
// The response body is unmarshalled into the Payload field of the Response
// object, using the rawdata.Unmarshal function, which supports JSON, YAML,
// TOML and other formats.
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
