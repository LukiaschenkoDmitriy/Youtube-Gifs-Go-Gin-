package command

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dmytrii/youtube-gifs-chat/fixture"
)

func FixtureCommand() {
	fixture := flag.NewFlagSet("fixture", flag.ExitOnError)

	fixtureName := fixture.String("fixture", "undefined", "fixture name")

	fixture.Parse(os.Args[2:])

	if fixtureName == nil {
		fmt.Println("Available fixtures: -fixture 'comments'")
		os.Exit(1)
	}

	if *fixtureName == "undefined" {
		fmt.Println("Available fixtures: -fixture 'comments'")
		os.Exit(1)
	}

	switch *fixtureName {
	case "comments":
		commentFixture(fixture)
	default:
		fmt.Println("Available fixtures: -fixture 'comments'")
		os.Exit(1)
	}
}

func commentFixture(fixtureSet *flag.FlagSet) {
	userId := fixtureSet.String("user", "", "UUID of the user the comments are created on behalf of (required)")
	count := fixtureSet.Int("count", 0, "Number of comments to generate (required)")
	videoId := fixtureSet.String("video", "fixture-video", "video_id the comments are attached to")
	nestedChance := fixtureSet.Float64("nested", 0.5, "Probability (0..1) that a comment is nested (a reply)")
	fixtureSet.Parse(os.Args[3:])

	if *userId == "" {
		log.Fatal("the -user flag (user UUID) is required")
	}
	if *count <= 0 {
		log.Fatal("the -count flag must be greater than 0")
	}

	fixture.RunCommentsFixture(userId, count, videoId, nestedChance)
}
