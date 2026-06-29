// Package httpx provides HTTP response types, status codes, and helpers.
// This package has no framework dependency — it works with net/http only.
package httpx

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Code is an application-level error code.
// The convention is: code / 1000 = HTTP status code.
type Code int

// Application error codes. By convention code/1000 is the HTTP status, except
// CodeSuccess which maps to 200.
const (
	CodeSuccess Code = 0

	CodeBadRequest    Code = 400000
	CodeInvalidParams Code = 400001
	CodeInvalidJSON   Code = 400002
	CodeInvalidFile   Code = 400003

	CodeUnauthorized Code = 401000
	CodeTokenMissing Code = 401001
	CodeTokenExpired Code = 401002
	CodeTokenInvalid Code = 401003

	CodeForbidden       Code = 403000
	CodeAccessDenied    Code = 403001
	CodeAccountDisabled Code = 403002

	CodeNotFound       Code = 404000
	CodeRouteNotFound  Code = 404001
	CodeRecordNotFound Code = 404002

	CodeConflict      Code = 409000
	CodeAlreadyExists Code = 409001
	CodeStateInvalid  Code = 409002

	CodeTooManyRequests Code = 429000

	CodeInternalError      Code = 500000
	CodeDBError            Code = 500001
	CodeCacheError         Code = 500002
	CodeThirdPartyError    Code = 500003
	CodeFileOperationError Code = 500004

	CodeServiceUnavailable Code = 503000
)

var codeMessages = map[Code]string{
	CodeSuccess:            "success",
	CodeBadRequest:         "bad request",
	CodeInvalidParams:      "invalid parameters",
	CodeInvalidJSON:        "invalid json format",
	CodeInvalidFile:        "invalid file",
	CodeUnauthorized:       "unauthorized",
	CodeTokenMissing:       "access token is missing",
	CodeTokenExpired:       "access token expired",
	CodeTokenInvalid:       "access token invalid",
	CodeForbidden:          "forbidden",
	CodeAccessDenied:       "access denied",
	CodeAccountDisabled:    "account disabled",
	CodeNotFound:           "resource not found",
	CodeRouteNotFound:      "route not found",
	CodeRecordNotFound:     "record not found",
	CodeConflict:           "conflict",
	CodeAlreadyExists:      "resource already exists",
	CodeStateInvalid:       "invalid resource state",
	CodeTooManyRequests:    "too many requests",
	CodeInternalError:      "internal server error",
	CodeDBError:            "database error",
	CodeCacheError:         "cache error",
	CodeThirdPartyError:    "external service error",
	CodeFileOperationError: "file operation error",
	CodeServiceUnavailable: "service unavailable",
}

// String returns the human-readable message for code c.
func (c Code) String() string {
	if msg, ok := codeMessages[c]; ok {
		return msg
	}
	return codeMessages[CodeInternalError]
}

// HTTPStatus returns the HTTP status code corresponding to c.
// CodeSuccess maps to 200; others use code / 1000.
func (c Code) HTTPStatus() int {
	if c == CodeSuccess {
		return http.StatusOK
	}
	return int(c / 1000)
}

// ─── Response shapes ──────────────────────────────────────────────────────────

type baseResponse struct {
	Code      Code   `json:"code"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

func newBase(reqID string, code Code, msg string) baseResponse {
	return baseResponse{
		Code:      code,
		Message:   msg,
		RequestID: reqID,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// SuccessResponse wraps a successful response with data.
type SuccessResponse[T any] struct {
	baseResponse
	Data *T `json:"data,omitempty"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	baseResponse
	ErrorStack []string `json:"error_stack,omitempty"`
	Details    []any    `json:"details,omitempty"`
}

func (e *ErrorResponse) Error() string { return e.Message }

// PageData wraps paginated results.
type PageData[T any] struct {
	Items    []T  `json:"items"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	Total    int  `json:"total"`
	HasMore  bool `json:"has_more"`
}

// PageResponse wraps a paginated response.
type PageResponse[T any] struct {
	baseResponse
	Data *PageData[T] `json:"data,omitempty"`
}

// ─── Constructors ─────────────────────────────────────────────────────────────

// NewSuccessResponse builds a 200 success response.
func NewSuccessResponse[T any](reqID string, data *T) SuccessResponse[T] {
	return SuccessResponse[T]{
		baseResponse: newBase(reqID, CodeSuccess, CodeSuccess.String()),
		Data:         data,
	}
}

// NewErrorResponse builds an error response from code and err.
// Optional details are attached to the Details field.
func NewErrorResponse(reqID string, code Code, err error, details ...any) *ErrorResponse {
	msg := code.String()
	var stack []string
	if err != nil {
		msg = err.Error()
		stack = strings.Split(fmt.Sprintf("%+v", err), "\n")
	}
	return &ErrorResponse{
		baseResponse: newBase(reqID, code, msg),
		ErrorStack:   stack,
		Details:      details,
	}
}

// NewErrorResponsef builds an error response from a format string.
func NewErrorResponsef(reqID string, code Code, format string, a ...any) *ErrorResponse {
	return NewErrorResponse(reqID, code, fmt.Errorf(format, a...))
}

// NewPageResponse builds a paginated response.
func NewPageResponse[T any](reqID string, items []T, total, page, pageSize int) PageResponse[T] {
	return PageResponse[T]{
		baseResponse: newBase(reqID, CodeSuccess, CodeSuccess.String()),
		Data: &PageData[T]{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			HasMore:  total > page*pageSize,
		},
	}
}

// NewPanicResponse builds an error response for a panic recovery.
func NewPanicResponse(reqID string, r any, stack []string) *ErrorResponse {
	return &ErrorResponse{
		baseResponse: newBase(reqID, CodeInternalError, fmt.Sprintf("panic: %v", r)),
		ErrorStack:   stack,
	}
}

// ─── HTTP error code mapping ──────────────────────────────────────────────────

// HTTPStatusToCode maps an HTTP status code to an application Code.
func HTTPStatusToCode(httpStatus int) Code {
	switch httpStatus {
	case http.StatusBadRequest, http.StatusMethodNotAllowed,
		http.StatusRequestTimeout, http.StatusRequestEntityTooLarge,
		http.StatusUnsupportedMediaType:
		return CodeInvalidParams
	case http.StatusUnauthorized:
		return CodeUnauthorized
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound, http.StatusGone:
		return CodeNotFound
	case http.StatusConflict:
		return CodeConflict
	case http.StatusPreconditionFailed:
		return CodeStateInvalid
	case http.StatusTooManyRequests:
		return CodeTooManyRequests
	case http.StatusServiceUnavailable:
		return CodeServiceUnavailable
	default:
		return CodeInternalError
	}
}
