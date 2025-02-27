package router

import (
	"context"
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Ctx struct {
	ctx context.Context
	w   http.ResponseWriter
	r   *http.Request
}

func NewCtx() *Ctx {
	return &Ctx{ctx: context.Background()}
}

func (c *Ctx) SetMiddleware(writer http.ResponseWriter, request *http.Request) {
	c.w = writer
	c.r = request
}

func (c *Ctx) Context() context.Context {
	if c.r != nil {
		return c.r.Context()
	}
	return c.ctx
}

func (c *Ctx) SetContext(ctx context.Context) {
	c.r = c.r.WithContext(ctx)
}

func (c *Ctx) WrapHTTPHandler(h http.HandlerFunc) error {
	h(c.w, c.r)
	return nil
}

func (c *Ctx) Span() trace.Span {
	return trace.SpanFromContext(c.Context())
}

func (c *Ctx) JSON(data any) error {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(http.StatusOK)

	c.Span().SetAttributes(attribute.Int("http.status_code", http.StatusOK))

	return json.NewEncoder(c.w).Encode(data)
}

func (c *Ctx) Error(err error, code ...int) {
	if err != nil {
		status := http.StatusInternalServerError
		if len(code) > 0 {
			status = code[0]
		}

		c.Span().SetAttributes(
			attribute.Int("http.status_code", status),
			attribute.String("http.error", err.Error()))

		http.Error(c.w, err.Error(), status)
	}
}

func (c *Ctx) Query(p string) URLParam {
	return URLParam(c.r.URL.Query().Get(p))
}

func (c *Ctx) Params(p string) string {
	return c.r.PathValue(p)
}

func (c *Ctx) ParseBody(v any) error {
	return json.NewDecoder(c.r.Body).Decode(&v)
}
