package core_server

import (
	"context"
	"fmt"
	"time"

	core_logger "github.com/daf32/url-shortener-fiber/internal/core/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type ServerConfig struct {
	Port            string
	BaseURL         string
	ShutdownTimeout time.Duration
}

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

type Registrar interface {
	RegisterRoutes(r fiber.Router)
}

type routeGroup struct {
	prefix     string
	registrars []Registrar
}

type HTTPServer struct {
	app        *fiber.App
	cfg        ServerConfig
	log        *core_logger.Logger
	middleware []fiber.Handler
	groups     []routeGroup
}

type Option func(*HTTPServer)

func WithMiddleware(m ...fiber.Handler) Option {
	return func(s *HTTPServer) {
		s.middleware = append(s.middleware, m...)
	}
}

func WithRoutes(prefix string, regs ...Registrar) Option {
	return func(s *HTTPServer) {
		s.groups = append(s.groups, routeGroup{prefix: prefix, registrars: regs})
	}
}

func NewHTTPServer(
	cfg ServerConfig,
	log *core_logger.Logger,
	opts ...Option,
) *HTTPServer {
	s := &HTTPServer{
		app: fiber.New(fiber.Config{
			StructValidator: &structValidator{validate: validator.New()},
		}),
		cfg: cfg,
		log: log,
	}

	for _, opt := range opts {
		opt(s)
	}

	for _, m := range s.middleware {
		s.app.Use(m)
	}

	for _, g := range s.groups {
		router := s.app.Group(g.prefix)

		for _, reg := range g.registrars {
			reg.RegisterRoutes(router)
		}
	}

	return s
}

func (s *HTTPServer) Run(ctx context.Context) error {
	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		if err := s.app.Listen(":" + s.cfg.Port); err != nil {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("start HTTP server: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		if err := s.app.ShutdownWithTimeout(s.cfg.ShutdownTimeout); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}
