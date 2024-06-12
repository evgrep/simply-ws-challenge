package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type IRequestHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) SetRequestHandler(pattern string, handler IRequestHandler) *Service {
	authCheckHandler := NewAuthCheckHandler(handler)
	http.Handle(pattern, http.HandlerFunc(authCheckHandler.Handle))
	return s
}

func (s *Service) Run() {

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Println("Starting serving new connections.")

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Server shutdown complete.")
}
