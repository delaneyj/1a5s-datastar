package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/delaneyj/toolbelt"
	"github.com/go-chi/chi/v5"
	"github.com/starfederation/1a4s-datastar/sql"
	"github.com/starfederation/1a4s-datastar/sql/zz"
	datastar "github.com/starfederation/datastar/code/go/sdk"
	"zombiezen.com/go/sqlite"
)

type PokemonBattle struct {
	UpvoteID   int64 `json:"upvoteId"`
	DownvoteID int64 `json:"downvoteId"`
}

func setupVoteRoutes(r chi.Router, db *toolbelt.Database) error {
	r.Route("/vote", func(voteRouter chi.Router) {

		randomBattle := func(tx *sqlite.Conn) (left, right *zz.PokemonModel, err error) {
			res, err := zz.OnceRandomPokemon(tx, 2)
			if err != nil {
				return nil, nil, fmt.Errorf("error getting random pokemon: %w", err)
			}
			left = sql.RandomResToPokemonModel(res[0])
			right = sql.RandomResToPokemonModel(res[1])
			return left, right, nil
		}

		voteRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var (
				left, right *zz.PokemonModel
			)
			if err := db.ReadTX(r.Context(), func(tx *sqlite.Conn) (err error) {
				left, right, err = randomBattle(tx)
				return err
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			VotePage(left, right).Render(r.Context(), w)
		})

		voteRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
			battle := &PokemonBattle{}
			if err := datastar.ReadSignals(r, battle); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			sse := datastar.NewSSE(w, r)

			now := time.Now()
			var left, right *zz.PokemonModel
			if err := db.WriteTX(r.Context(), func(tx *sqlite.Conn) (err error) {
				upvoteErr := zz.OnceUpvotePokemon(tx, zz.UpvotePokemonParams{
					Id:        battle.UpvoteID,
					UpdatedAt: now,
				})
				downvoteErr := zz.OnceDownvotePokemon(tx, zz.DownvotePokemonParams{
					Id:        battle.DownvoteID,
					UpdatedAt: now,
				})
				if err := errors.Join(upvoteErr, downvoteErr); err != nil {
					return fmt.Errorf("error voting: %w", err)
				}

				// Get new battle
				left, right, err = randomBattle(tx)
				return err
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			sse.MergeFragmentTempl(voteContainer(left, right))
		})

	})
	return nil
}
