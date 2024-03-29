package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenAIErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type OpenAIClient struct {
	config Configuration
	client *http.Client
}

func NewClient(config Configuration) (openai *OpenAIClient, err error) {
	client := http.DefaultClient
	openai = &OpenAIClient{config, client}
	return
}

func (openai OpenAIClient) MakeRequest(path string, data interface{}) (io.ReadCloser, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json error: %v", err)
	}
	req, err := http.NewRequest("POST", openai.config.API+path, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("invalid request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+openai.config.APIKey)
	res, err := openai.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot make request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %s", res.Status)
	}
	return res.Body, nil
}
