package internalhttp

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net"
	"net/http"
	"rotator/internal/app"
)

type Server struct {
	host   string
	port   string
	logger Logger
	server *http.Server
}

type Logger interface {
	Debug(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	LogHTTP(r *http.Request, code, length int)
}

func NewServer(host, port string, app *app.App, logger Logger) *Server {
	server := &Server{
		host:   host,
		port:   port,
		logger: logger,
		server: nil,
	}

	httpServ := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: loggingMiddleware(Routers(app), logger),
	}

	server.server = httpServ

	return server
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("[+] Staring http server and listen", zap.String("host:", s.host), zap.String("port", s.port))
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func Routers(app *app.App) http.Handler {
	handlers := NewServerHandlers(app)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/banner-slot/add", handlers.AddBannerToSlot).Methods("POST")
	r.HandleFunc("/api/v1/banner-slot/remove", handlers.RemoveBannerToSlot).Methods("DELETE")
	r.HandleFunc("/api/v1/banner/transition", handlers.CountTransition).Methods("POST")
	r.HandleFunc("/api/v1/banner/choose", handlers.ChooseBanner).Methods("POST")

	return r
}
