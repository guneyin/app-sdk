package router

import (
	"fmt"
	"net/http"

	"github.com/guneyin/app-sdk/logger"
	"go.opentelemetry.io/otel"
)

type Router struct {
	logger *logger.Logger
	mux    *http.ServeMux
}

func newRouter(logger *logger.Logger, mux *http.ServeMux) *Router {
	return &Router{logger, mux}
}

func (r *Router) registerHandler(method HTTPMethod, path string, handler HandlerFunc) {
	logger.Warn("%-7s | %-50s\n", method, path)

	r.mux.HandleFunc(fmt.Sprintf("%s %s", method, path), func(writer http.ResponseWriter, request *http.Request) {
		ctx := NewCtx()
		ctx.SetMiddleware(writer, request)

		tracer := otel.Tracer(fmt.Sprintf("%s %s", method, path))
		spanCtx, span := tracer.Start(ctx.Context(), fmt.Sprintf("%s %s", method, path))
		defer span.End()

		ctx.SetContext(spanCtx)

		writer.Header().Set("TraceID", span.SpanContext().TraceID().String())

		defer func(c *Ctx) {
			if rec := recover(); rec != nil {
				errPanic := fmt.Errorf("panic recovered: %v", rec)

				if c != nil {
					c.Error(errPanic)
				} else {
					r.logger.Error(errPanic.Error())
				}
			}
		}(ctx)

		err := handler(ctx)
		if err != nil {
			ctx.Error(err)
		}
	})
}
