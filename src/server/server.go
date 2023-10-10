package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port   = ":8090"
	errGrp errgroup.Group
)

type Server struct {
	Router chi.Router
}

func NewServer() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	srv := &Server{
		Router: r,
	}
	srv.registerRoutes()
	return r
}
func Start(r *chi.Mux) {

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	errGrp.Go(func() error {
		log.Info().Msgf("starting server in port: %v", port)
		srv := &http.Server{
			Addr:    port,
			Handler: r,
		}

		// Listen for syscall signals for process to interrupt/quit
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			<-sig

			log.Info().Msg("received interrupt")
			// Shutdown signal with grace period of 30 seconds
			shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

			go func() {
				<-shutdownCtx.Done()
				if shutdownCtx.Err() == context.DeadlineExceeded {
					log.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
				}
			}()

			// Trigger graceful shutdown
			err := srv.Shutdown(shutdownCtx)
			if err != nil {
				log.Fatal().Err(err)
			}
			serverStopCtx()
		}()

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err)
		}
		return err
	})

	if err := errGrp.Wait(); err != nil {
		log.Fatal().Err(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

}

func (srv *Server) registerRoutes() {
	srv.Router.Get("/health", health)
}

func health(resp http.ResponseWriter, r *http.Request) {
	resp.WriteHeader(http.StatusOK)
}
