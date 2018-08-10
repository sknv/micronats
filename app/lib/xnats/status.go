package xnats

import (
	"fmt"
	"net/http"
)

type StatusCode uint32

const (
	StatusOK               StatusCode = 0
	StatusInvalidArgument  StatusCode = 1
	StatusUnauthenticated  StatusCode = 2
	StatusPermissionDenied StatusCode = 3
	StatusInternal         StatusCode = 4
	StatusDeadlineExceeded StatusCode = 5
)

func ServerHTTPStatusFromErrorCode(code StatusCode) int {
	switch code {
	case StatusOK:
		return http.StatusOK
	case StatusInvalidArgument:
		return http.StatusBadRequest
	case StatusUnauthenticated:
		return http.StatusUnauthorized
	case StatusPermissionDenied:
		return http.StatusForbidden
	case StatusInternal:
		return http.StatusInternalServerError
	case StatusDeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return 0 // invalid
	}
}

func IsValidErrorCode(code StatusCode) bool {
	return ServerHTTPStatusFromErrorCode(code) != 0
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type Status struct {
	Code    StatusCode        `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Meta    map[string]string `json:"meta,omitempty"`
}

func ErrorStatus(code StatusCode, message string) *Status {
	if IsValidErrorCode(code) {
		return &Status{
			Code:    code,
			Message: message,
		}
	}
	return &Status{
		Code:    StatusInternal,
		Message: "rpc: invalid status code: " + fmt.Sprint(code),
	}
}

func StatusFromError(err error) (*Status, bool) {
	status, match := err.(*Status)
	if match {
		return status, true
	}
	return ErrorStatus(StatusInternal, err.Error()), false
}

func (s *Status) StatusCode() StatusCode {
	if s == nil {
		return StatusOK
	}
	return StatusCode(s.Code)
}

func (s *Status) MetaValue(key string) string {
	if s.Meta != nil {
		return s.Meta[key] // also returns "" if key is not in meta map
	}
	return ""
}

func (s *Status) WithMeta(key string, value string) {
	if s.Meta == nil {
		s.Meta = make(map[string]string)
	}
	s.Meta[key] = value // upsert the value
}

func (s *Status) Error() string {
	return fmt.Sprintf("rpc error %d: %s", s.Code, s.Message)
}
