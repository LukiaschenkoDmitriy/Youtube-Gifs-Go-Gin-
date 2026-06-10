package main

import (
	"embed"
	"log"

	"github.com/dmytrii/youtube-gifs-chat/config"
	"github.com/dmytrii/youtube-gifs-chat/internal/usecase"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func main() {
	c := config.GetConfig()
	pool, err := usecase.NewDatabase(c)

	if err != nil {
		log.Fatal(err)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	goose.SetDialect("postgres")
	goose.SetBaseFS(MigrationsFS)

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal("migration failed: %w", err)
	}
}
