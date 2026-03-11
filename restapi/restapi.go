package restapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type ApiPack struct {
	Url    string `json:"url"`
	Method string `json:"method"`
	Token  string `json:"token,omitempty"`
	Body   any    `json:"body,omitempty"`
}

func Call(api *ApiPack) ([]byte, error) {
	// body
	var bodyReader io.Reader
	if api.Body != nil {
		jsonBody, err := json.Marshal(api.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}
	// request
	req, err := http.NewRequest(api.Method, api.Url, bodyReader)
	if err != nil {
		return nil, err
	}
	// header
	if api.Token != "" {
		req.Header.Set("Authorization", "Bearer "+api.Token)
	}
	if api.Body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	// client
	client := &http.Client{Timeout: time.Second * 10}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	// response
	if rsp.StatusCode >= 400 {
		return nil, errors.New(rsp.Status)
	}
	return io.ReadAll(rsp.Body)
}
