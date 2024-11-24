package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/delaneyj/toolbelt"
	"github.com/go-chi/chi/v5"
	"github.com/starfederation/1a4s-datastar/sql/zz"
	datastar "github.com/starfederation/datastar/code/go/sdk"
	"zombiezen.com/go/sqlite"
)

const blockSize = 64

func setupResultsRoutes(r chi.Router, db *toolbelt.Database) error {
	rows := []zz.ResultsRes{}
	mu := &sync.RWMutex{}

	updateResults := toolbelt.Throttle(5*time.Second, func(ctx context.Context) error {
		mu.Lock()
		defer mu.Unlock()

		rows = rows[:0]
		if err := db.ReadTX(ctx, func(tx *sqlite.Conn) (err error) {
			rows, err = zz.OnceResults(tx)
			return err
		}); err != nil {
			return fmt.Errorf("error reading results: %w", err)
		}

		return nil
	})

	r.Route("/results", func(resultsRouter chi.Router) {
		resultsRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			updateResults(ctx)
			ResultsPage().Render(ctx, w)
		})

		resultsRouter.Get("/{blockOffsetRaw}", func(w http.ResponseWriter, r *http.Request) {
			blockOffsetRaw := chi.URLParam(r, "blockOffsetRaw")
			blockOffset, err := strconv.Atoi(blockOffsetRaw)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			sse := datastar.NewSSE(w, r)

			mu.RLock()
			defer mu.RUnlock()
			blockRows := rows[blockOffset*blockSize : (blockOffset+1)*blockSize]
			sse.MergeFragmentTempl(
				resultsRows(blockOffset, blockRows...),
				datastar.WithSelectorID("results"),
				datastar.WithMergeAppend(),
			)

			sse.MergeFragmentTempl(loadNextBlock(blockOffset))
		})
	})
	return nil
}
