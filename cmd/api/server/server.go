package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Addr   string
}

func NewServer(router *mux.Router, addr string) *Server {
	return &Server{
		Router: router,
		Addr:   addr,
	}
}

func (s *Server) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:    s.Addr,
		Handler: s.Router,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
		ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxShutDown); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
		}
	}()

	log.Printf("Server running on %s", s.Addr)
	return srv.ListenAndServe()
}
