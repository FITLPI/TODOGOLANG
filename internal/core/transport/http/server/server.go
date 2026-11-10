package core_http_server

import (
	"context"
	"fmt"
	"net/http"

	core_logger "github.com/FITLPI/TODOGOLANG/internal/core/logger"
	core_http_middleware "github.com/FITLPI/TODOGOLANG/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	middlewares []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middlewares ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:         http.NewServeMux(),
		config:      config,
		log:         log,
		middlewares: middlewares,
	}
}

func (h *HTTPServer) RegistrAPIRoutes(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := fmt.Sprintf("/api/%s", router.apiVersion)
		h.mux.Handle(prefix+"/", http.StripPrefix(prefix, router))
	}
}

func (h *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMidlewares(h.mux, h.middlewares...)

	server := &http.Server{
		Addr:    h.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)

		h.log.Warn("starting server", zap.String("addr", h.config.Addr))
		err := server.ListenAndServe()

		if err != http.ErrServerClosed {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		return fmt.Errorf("listen and serv HTTP: %w", err)
	case <-ctx.Done():
		h.log.Warn("shutting down server", zap.String("addr", h.config.Addr))
		shutdownCtx, cancel := context.WithTimeout(context.Background(), h.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		h.log.Warn("server shutdown complete", zap.String("addr", h.config.Addr))

	}

	return nil
}
