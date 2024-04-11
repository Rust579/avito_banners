package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
)

type Error struct {
	Key int    `json:"key"`
	Msg string `json:"message"`
}

type Response struct {
	Status    bool          `json:"status"`
	Errors    []Error       `json:"errors"`
	Values    []interface{} `json:"values"`
	TmRequest string        `json:"tm_req"`
}

// SendHTTPRequest http клиент для тестов
func SendHTTPRequest(uri string, method string, headers map[string]string, statusCode int, reqBody interface{}) (*Response, error) {
	client := &fasthttp.Client{}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI("http://localhost:8000" + uri)
	req.Header.SetMethod(method)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)

	err = client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != statusCode {
		fmt.Println(string(resp.Body()))
		return nil, errors.New("неправильный статус-код")
	}

	var response Response

	if resp.Body() != nil {
		err = json.Unmarshal(resp.Body(), &response)
		if err != nil {
			return nil, err
		}
	}

	return &response, nil
}
