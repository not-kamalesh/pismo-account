package dto

import (
	"log/slog"
	"net/http"
)

type RequestParserValidator interface {
	Parse(r *http.Request) error
	Validate() error
}

func ParseRequest[T RequestParserValidator](r *http.Request, newReq func() T) (T, error) {
	req := newReq()
	if err := req.Parse(r); err != nil {
		slog.Warn("failed to parse request", "type", req, "error", err)
		return req, err
	}
	if err := req.Validate(); err != nil {
		slog.Warn("validation failed", "type", req, "error", err)
		return req, err
	}
	return req, nil
}
