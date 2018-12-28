package nats

import (
	"fmt"
	"net/http"
)

type StatusCode int32

const (
	OK               StatusCode = 0
	InvalidArgument  StatusCode = 1
	Unauthenticated  StatusCode = 2
	PermissionDenied StatusCode = 3
	Internal         StatusCode = 4
	DeadlineExceeded StatusCode = 5
)

func ServerHTTPStatusFromErrorCode(code StatusCode) int {
	switch code {
	case OK:
		return http.StatusOK
	case InvalidArgument:
		return http.StatusBadRequest
	case Unauthenticated:
		return http.StatusUnauthorized
	case PermissionDenied:
		return http.StatusForbidden
	case Internal:
		return http.StatusInternalServerError
	case DeadlineExceeded:
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

func StatusError(code StatusCode, message string) *Status {
	if IsValidErrorCode(code) {
		return &Status{
			Code:    int32(code),
			Message: message,
		}
	}
	return &Status{
		Code:    int32(Internal),
		Message: "invalid status code: " + fmt.Sprint(code),
	}
}

func StatusFromError(err error) (*Status, bool) {
	status, match := err.(*Status)
	if match {
		return status, true
	}
	return StatusError(Internal, err.Error()), false
}

func (s *Status) StatusCode() StatusCode {
	if s == nil {
		return OK
	}
	return StatusCode(s.Code)
}

func (s *Status) HasError() bool {
	return s.StatusCode() != OK
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
	return fmt.Sprintf("%d: %s", s.Code, s.Message)
}
