package main

import (
	"embed"
	"log"
	"os"

	"github.com/dmytrii/youtube-gifs-chat/config"
	"github.com/dmytrii/youtube-gifs-chat/internal/database"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	if len(os.Args) < 2 || (os.Args[1] != "up" && os.Args[1] != "down") {
		log.Fatal("usage: migrate <up|down>")
	}

	c, err := config.GetConfig()

	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.NewDatabase(c)

	if err != nil {
		log.Fatal(err)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	goose.SetBaseFS(migrationsFS)

	if os.Args[1] == "up" {
		err = goose.Up(db, "migrations")
	} else {
		err = goose.Down(db, "migrations")
	}

	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
