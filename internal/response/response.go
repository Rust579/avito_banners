package response

import (
	"avito_banners/internal/errs"
	"encoding/json"
	"reflect"
	"time"
)

type Response struct {
	Status      bool          `json:"status"`
	Errors      []errs.Error  `json:"errors"`
	Values      []interface{} `json:"values"`
	TmRequest   string        `json:"tm_req"`
	TmRequestSt time.Time     `json:"-"`
}

func InitResponse() *Response {
	return &Response{
		TmRequestSt: time.Now(),
	}
}

func (r *Response) SetError(err errs.Error) *Response {
	r.Errors = append(r.Errors, err)
	return r
}

func (r *Response) SetErrors(ers []errs.Error) *Response {
	r.Errors = append(r.Errors, ers...)

	return r
}

func (r *Response) SetValue(val interface{}) *Response {
	r.Values = append(r.Values, val)
	return r
}

func (r *Response) SetValues(vals interface{}) *Response {
	r.Values = r.InterfaceSlice(vals)
	return r
}

func (r *Response) InterfaceSlice(in interface{}) []interface{} {
	s := reflect.ValueOf(in)
	if s.Kind() != reflect.Slice {
		return []interface{}{}
	}
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}

func (r *Response) FormResponse() *Response {
	if len(r.Errors) > 0 {
		r.Status = false
	} else {
		r.Status = true
	}

	// not null array json
	if len(r.Values) == 0 {
		r.Values = []interface{}{}
	}

	// not null array json
	if len(r.Errors) == 0 {
		r.Errors = []errs.Error{}
	}
	return r
}

func (r *Response) Json() []byte {
	r.TmRequest = time.Now().Sub(r.TmRequestSt).String()
	if bts, err := json.Marshal(r); err == nil {
		return bts
	}
	return []byte{}
}
