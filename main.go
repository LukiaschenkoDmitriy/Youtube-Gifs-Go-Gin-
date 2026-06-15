package main

import (
	"fmt"
	"os"

	"github.com/dmytrii/youtube-gifs-chat/command"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Commands:")
		fmt.Println("server - start GIN server")
		fmt.Println("fixture -fixture=[entity] - start fixture")
		fmt.Println("migrate [up|down] - migrate .sql files")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		command.ServerCommand()
	case "fixture":
		command.FixtureCommand()
	case "migrate":
		command.MigrateCommand()
	default:
		fmt.Printf("Unknown command: %q\n", os.Args[1])
		os.Exit(1)
	}

}
