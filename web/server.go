package web

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/benbjohnson/hashfs"
	"github.com/delaneyj/toolbelt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static/*
var staticFS embed.FS

var staticSys = hashfs.NewFS(staticFS)

func staticPath(path string) string {
	return "/" + staticSys.HashName("static/"+path)
}

func RunBlocking(db *toolbelt.Database, port int) toolbelt.CtxErrFunc {
	return func(ctx context.Context) error {
		router := chi.NewRouter()
		router.Use(middleware.Recoverer)
		router.Handle("/static/*", hashfs.FileServer(staticSys))

		if err := errors.Join(
			setupHomeRoutes(router, db),
			setupVoteRoutes(router, db),
			setupResultsRoutes(router, db),
		); err != nil {
			return fmt.Errorf("error setting up routes: %w", err)
		}

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		}
		go func() {
			<-ctx.Done()
			srv.Shutdown(context.Background())
		}()

		log.Printf("Hosting on http://localhost:%d", port)
		return srv.ListenAndServe()
	}
}
