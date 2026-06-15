package command

import (
	"fmt"
	"os"

	"github.com/dmytrii/youtube-gifs-chat/migrate"
)

func MigrateCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Available arguments: 'up' | 'down'")
		os.Exit(1)
	}

	migrateArg := os.Args[2]

	if migrateArg != "down" && migrateArg != "up" {
		fmt.Println("Available arguments: 'up' | 'down'")
		os.Exit(1)
	}

	migrate.RunMigration(migrateArg)
}
