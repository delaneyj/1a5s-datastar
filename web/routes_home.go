package web

import (
	"fmt"
	"net/http"

	"github.com/delaneyj/toolbelt"
	"github.com/go-chi/chi/v5"
	"github.com/starfederation/1a4s-datastar/sql/zz"
	datastar "github.com/starfederation/datastar/code/go/sdk"
	"zombiezen.com/go/sqlite"
)

func setupHomeRoutes(r chi.Router, db *toolbelt.Database) error {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/vote", http.StatusSeeOther)
	})

	r.Get("/prefetch", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)

		var ids []int64
		if err := db.ReadTX(r.Context(), func(tx *sqlite.Conn) (err error) {
			ids, err = zz.OnceAllIds(tx)
			return err
		}); err != nil {
			sse.ConsoleError(err)
			return
		}

		urls := make([]string, len(ids))
		for i, id := range ids {
			urls[i] = fmt.Sprintf(pokemonSpriteURLFormat, id)
		}
		sse.Prefetch(urls...)
	})

	return nil
}
