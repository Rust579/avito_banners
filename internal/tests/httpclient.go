package tests

import (
	"encoding/json"
	"errors"
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

func SendHTTPRequest(uri string, method string, headers map[string]string, reqBody interface{}) (*Response, error) {
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

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, errors.New("неправильный статус-код")
	}

	var response Response

	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
