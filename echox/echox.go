// Package echox bridges gokit's framework-agnostic httpx and contextx packages
// into the labstack/echo web framework.
//
// It provides response helpers that emit httpx's standard envelope, a
// RequestID middleware that propagates a request ID through the context, and an
// echo.HTTPErrorHandler that renders errors using the httpx schema.
package echox

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/tmoeish/gokit/contextx"
	"github.com/tmoeish/gokit/httpx"
	"github.com/tmoeish/gokit/randx"
)

// Success writes a 200 response wrapping data in the httpx success envelope.
func Success[T any](c echo.Context, data *T) error {
	reqID := contextx.ReqID(c.Request().Context())
	return c.JSON(http.StatusOK, httpx.NewSuccessResponse(reqID, data))
}

// NoContent writes a 200 success response with no data payload.
func NoContent(c echo.Context) error {
	return Success[struct{}](c, nil)
}

// Error writes an error response using code's HTTP status and message.
func Error(c echo.Context, code httpx.Code, err error, details ...any) error {
	reqID := contextx.ReqID(c.Request().Context())
	resp := httpx.NewErrorResponse(reqID, code, err, details...)
	return c.JSON(code.HTTPStatus(), resp)
}

// Page writes a paginated 200 response using the httpx page envelope.
func Page[T any](c echo.Context, items []T, total, page, pageSize int) error {
	reqID := contextx.ReqID(c.Request().Context())
	return c.JSON(http.StatusOK, httpx.NewPageResponse(reqID, items, total, page, pageSize))
}

// RequestID returns middleware that ensures every request carries a request ID.
//
// It reuses an inbound X-Request-ID header when present, otherwise generates a
// UUID. The ID is stored in the request context via contextx.WithReqID (so it
// flows into logs and downstream handlers) and echoed back in the response
// header.
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = randx.UUID()
			}
			ctx := contextx.WithReqID(req.Context(), rid)
			c.SetRequest(req.WithContext(ctx))
			c.Response().Header().Set(echo.HeaderXRequestID, rid)
			return next(c)
		}
	}
}

// ErrorHandler is an echo.HTTPErrorHandler that renders errors using the httpx
// envelope. Register it with e.HTTPErrorHandler = echox.ErrorHandler.
//
// *echo.HTTPError values are mapped to the closest httpx.Code; *httpx.ErrorResponse
// values are rendered as-is; any other error becomes a CodeInternalError.
func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	reqID := contextx.ReqID(c.Request().Context())

	switch e := err.(type) {
	case *httpx.ErrorResponse:
		_ = c.JSON(e.Code.HTTPStatus(), e)
	case *echo.HTTPError:
		code := httpx.HTTPStatusToCode(e.Code)
		resp := httpx.NewErrorResponsef(reqID, code, "%v", e.Message)
		_ = c.JSON(code.HTTPStatus(), resp)
	default:
		resp := httpx.NewErrorResponse(reqID, httpx.CodeInternalError, err)
		_ = c.JSON(http.StatusInternalServerError, resp)
	}
}
