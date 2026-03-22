package dto

import (
	"net/http"
)

type RequestParserValidator interface {
	Parse(r *http.Request) error
	Validate() error
}

func ParseRequest[T RequestParserValidator](r *http.Request, newReq func() T) (T, error) {
	req := newReq()
	if err := req.Parse(r); err != nil {
		return req, err
	}
	if err := req.Validate(); err != nil {
		return req, err
	}
	return req, nil
}
