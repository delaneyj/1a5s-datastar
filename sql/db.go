package sql

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"embed"

	"github.com/delaneyj/toolbelt"
	"github.com/goccy/go-json"
	"github.com/starfederation/1a4s-datastar/sql/zz"
	"github.com/valyala/bytebufferpool"
	"zombiezen.com/go/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func New(ctx context.Context) (*toolbelt.Database, error) {
	migrations, err := toolbelt.MigrationsFromFS(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("err creating migrations: %w", err)
	}

	db, err := toolbelt.NewDatabase(ctx, "pokemon.sqlite", migrations)
	if err != nil {
		return nil, fmt.Errorf("err creating database: %w", err)
	}

	var isEmpty bool
	if err := db.ReadTX(ctx, func(tx *sqlite.Conn) error {
		count, err := zz.OnceCountPokemon(tx)
		if err != nil {
			return fmt.Errorf("error counting pokemon: %w", err)
		}
		isEmpty = count == 0
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error checking if database is empty: %w", err)
	}

	if isEmpty {
		log.Print("Empty database, seeding")
		if err := Seed(db); err != nil {
			return nil, fmt.Errorf("error seeding database: %w", err)
		}
	}

	return db, nil
}

func Seed(db *toolbelt.Database) error {
	const q = `
query GetAllPokemon {
	pokemon_v2_pokemon {
		id
		pokemon_v2_pokemonspecy {
			name
		}
	}
}`
	b, err := json.Marshal(map[string]string{"query": q})
	if err != nil {
		return fmt.Errorf("error marshalling query: %w", err)
	}
	r := bytes.NewReader(b)

	req, err := http.NewRequest("POST", "https://beta.pokeapi.co/graphql/v1beta", r)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	if _, err := buf.ReadFrom(res.Body); err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error response code: %d", res.StatusCode)
	}

	type GraphQLResponse struct {
		Data struct {
			Pokemon []struct {
				ID           int64 `json:"id"`
				PokemonSpecy struct {
					Name string `json:"name"`
				} `json:"pokemon_v2_pokemonspecy"`
			} `json:"pokemon_v2_pokemon"`
		} `json:"data"`
	}
	gqlRes := &GraphQLResponse{}
	if err := json.Unmarshal(buf.Bytes(), gqlRes); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	if err := db.WriteTX(context.Background(), func(tx *sqlite.Conn) error {
		insert := zz.CreatePokemon(tx)
		now := time.Now()
		for _, p := range gqlRes.Data.Pokemon {
			// https://github.com/t3dotgg/1app5stacks/blob/main/go-graphql-spa-version/go-gql-server/main.go#L144C108-L144C124
			if p.ID >= 1025 {
				continue
			}
			if err := insert.Run(&zz.PokemonModel{
				Id:         p.ID,
				Name:       p.PokemonSpecy.Name,
				DexId:      p.ID,
				InsertedAt: now,
				UpdatedAt:  now,
			}); err != nil {
				return fmt.Errorf("error inserting pokemon: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("error inserting pokemon: %w", err)
	}

	return nil
}

func RandomResToPokemonModel(p zz.RandomPokemonRes) *zz.PokemonModel {
	return &zz.PokemonModel{
		Id:         p.Id,
		Name:       p.Name,
		DexId:      p.DexId,
		UpVotes:    p.UpVotes,
		DownVotes:  p.DownVotes,
		InsertedAt: p.InsertedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}
