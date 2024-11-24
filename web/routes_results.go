package web

import (
	"net/http"

	"github.com/delaneyj/toolbelt"
	"github.com/go-chi/chi/v5"
	"github.com/starfederation/1a4s-datastar/sql/zz"
	datastar "github.com/starfederation/datastar/code/go/sdk"
	"zombiezen.com/go/sqlite"
)

func setupResultsRoutes(r chi.Router, db *toolbelt.Database) error {
	r.Route("/results", func(resultsRouter chi.Router) {
		resultsRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			ResultsPage().Render(r.Context(), w)
		})

		resultsRouter.Get("/rows", func(w http.ResponseWriter, r *http.Request) {
			sse := datastar.NewSSE(w, r)
			ctx := r.Context()
			var rows []zz.ResultsRes
			if err := db.ReadTX(ctx, func(tx *sqlite.Conn) (err error) {
				rows, err = zz.OnceResults(tx)
				return err
			}); err != nil {
				sse.ConsoleError(err)
				return
			}
			sse.MergeFragmentTempl(resultRows(rows...))
		})
	})
	return nil
}
