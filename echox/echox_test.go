package echox_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/tmoeish/gokit/echox"
	"github.com/tmoeish/gokit/httpx"
)

func newCtx(e *echo.Echo, method, target string) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestSuccess(t *testing.T) {
	e := echo.New()
	c, rec := newCtx(e, http.MethodGet, "/")
	type payload struct {
		Name string `json:"name"`
	}
	if err := echox.Success(c, &payload{Name: "gokit"}); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"name":"gokit"`) {
		t.Fatalf("body = %s", rec.Body.String())
	}
}

func TestError(t *testing.T) {
	e := echo.New()
	c, rec := newCtx(e, http.MethodGet, "/")
	if err := echox.Error(c, httpx.CodeNotFound, errors.New("missing")); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	e := echo.New()
	c, rec := newCtx(e, http.MethodGet, "/")

	var seenReqID string
	handler := echox.RequestID()(func(c echo.Context) error {
		seenReqID = c.Request().Header.Get(echo.HeaderXRequestID)
		return echox.NoContent(c)
	})
	if err := handler(c); err != nil {
		t.Fatal(err)
	}
	if got := rec.Header().Get(echo.HeaderXRequestID); got == "" {
		t.Fatal("response missing X-Request-ID header")
	}
	_ = seenReqID
}

func TestErrorHandler(t *testing.T) {
	e := echo.New()
	c, rec := newCtx(e, http.MethodGet, "/")
	echox.ErrorHandler(echo.NewHTTPError(http.StatusBadRequest, "bad"), c)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}
