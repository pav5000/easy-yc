package functions

import (
	"encoding/base64"
	"encoding/json"
)

type Request struct {
	HttpMethod      string `json:"httpMethod"`
	Body            string `json:"body"`
	IsBase64Encoded bool   `json:"isBase64Encoded"`
}

func (r *Request) UnmarshalBody(v interface{}) error {
	var rawJSON []byte
	if r.IsBase64Encoded {
		var err error
		rawJSON, err = base64.StdEncoding.DecodeString(r.Body)
		if err != nil {
			return err
		}
	} else {
		rawJSON = []byte(r.Body)
	}
	return json.Unmarshal(rawJSON, v)
}

type Response struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}
