package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/yash492/statusy/pkg/types"
)

type result struct {
	statusCode int
	data       any
	meta       types.JSON
}

func (r *result) New(statusCode int, data any, meta types.JSON) *result {
	return &result{
		statusCode: statusCode,
		data:       data,
		meta:       meta,
	}
}

func (r *result) sendHelper(w http.ResponseWriter) *Response {
	return &Response{
		Status: r.statusCode,
		Data: types.JSON{
			"data": r.data,
			"meta": r.meta,
		},
		IsError: false,
	}
}

func Send(w http.ResponseWriter, statusCode int, data any, meta types.JSON) *Response {
	r := &result{}
	return r.New(statusCode, data, meta).sendHelper(w)
}

func Errorf(w http.ResponseWriter, statusCode int, msg string, args ...any) *Response {
	r := &result{}
	return r.New(statusCode, fmt.Sprintf(msg, args...), nil).sendErrorHelper(w)
}

func (r *result) sendErrorHelper(w http.ResponseWriter) *Response {
	return &Response{
		Data: types.JSON{
			"is_error":    true,
			"error_msg":   r.data,
			"status_code": r.statusCode,
		},
		Status:  r.statusCode,
		IsError: true,
	}
}

func MiddlewareErrorf(w http.ResponseWriter, statusCode int, msg string, args ...any) {
	res := Response{
		Data: types.JSON{
			"is_error":    true,
			"error_msg":   fmt.Sprintf(msg, args...),
			"status_code": statusCode,
		},
		Status:  statusCode,
		IsError: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Status)

	if res.Status == http.StatusNoContent {
		return
	}

	err := json.NewEncoder(w).Encode(res.Data)
	if err != nil {
		slog.Error(err.Error())
	}
}
