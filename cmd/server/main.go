package main

import (
	"log"

	"github.com/daf32/url-shortener-fiber/internal/logger"
	"github.com/daf32/url-shortener-fiber/internal/repository"
	"github.com/daf32/url-shortener-fiber/internal/service"
	transport "github.com/daf32/url-shortener-fiber/internal/transport/http"
)

func main() {
	repo := repository.NewMemoryRepo()
	svc := service.NewShortenerService(repo)
	handler := transport.NewHTTPHandler(svc)

	logWriter, err := logger.NewLogWriter("logs")
	if err != nil {
		log.Fatal(err)
	}
	app := transport.NewRouter(handler, logWriter)

	log.Fatal(app.Listen(":3000"))
}
