package web

import (
	"fmt"
	"net/http"

	"github.com/delaneyj/toolbelt"
	"github.com/go-chi/chi/v5"
	"github.com/starfederation/1a4s-datastar/sql/zz"
	"zombiezen.com/go/sqlite"
)

func setupResultsRoutes(r chi.Router, db *toolbelt.Database) error {
	r.Route("/results", func(resultsRouter chi.Router) {
		resultsRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var rows []zz.ResultsRes
			if err := db.ReadTX(ctx, func(tx *sqlite.Conn) (err error) {
				rows, err = zz.OnceResults(tx)
				return err
			}); err != nil {
				http.Error(w, fmt.Sprintf("error reading results: %v", err), http.StatusInternalServerError)
				return
			}
			ResultsPage(rows...).Render(ctx, w)
		})
	})
	return nil
}
