package main

import (
	"context"
	"fmt"
	"log"

	"github.com/starfederation/1a4s-datastar/sql"
	"github.com/starfederation/1a4s-datastar/web"
)

const port = 4321

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	db, err := sql.New(ctx)
	if err != nil {
		return fmt.Errorf("error creating database: %w", err)
	}
	defer db.Close()

	return web.RunBlocking(db, port)(ctx)

}
