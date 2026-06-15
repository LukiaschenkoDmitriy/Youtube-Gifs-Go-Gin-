package fixture

import (
	"context"
	"log"
	"math/rand"

	"github.com/dmytrii/youtube-gifs-chat/config"
	"github.com/dmytrii/youtube-gifs-chat/internal/database"
)

// gifUrls — list of gif links.
// TODO: fill in with the links you will provide later.
var gifUrls = []string{
	"https://media4.giphy.com/media/v1.Y2lkPTc5MGI3NjExYXgwbG5hNGhsN3prZDl3enRldHhjcHRoeWYxbzVmYmpqbXA5d25jdiZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/XuGI90AtuG8PhWKZBX/giphy.gif",
	"https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExM3B2MmtvdW55YTRndnFyY2g5dGpxMzlmbmZucnR2czViOHphaGtpeCZlcD12MV9naWZzX3RyZW5kaW5nJmN0PWc/XIU0gdY0OYnYYz6GC3/giphy.gif",
	"https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExM3B2MmtvdW55YTRndnFyY2g5dGpxMzlmbmZucnR2czViOHphaGtpeCZlcD12MV9naWZzX3RyZW5kaW5nJmN0PWc/8IPKZHrNpSTO8/giphy.gif",
	"https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExM3B2MmtvdW55YTRndnFyY2g5dGpxMzlmbmZucnR2czViOHphaGtpeCZlcD12MV9naWZzX3RyZW5kaW5nJmN0PWc/XHeLeuirRbwptHhSWd/giphy.gif",
	"https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExM3B2MmtvdW55YTRndnFyY2g5dGpxMzlmbmZucnR2czViOHphaGtpeCZlcD12MV9naWZzX3RyZW5kaW5nJmN0PWc/5VKbvrjxpVJCM/giphy.gif",
	"https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExM3B2MmtvdW55YTRndnFyY2g5dGpxMzlmbmZucnR2czViOHphaGtpeCZlcD12MV9naWZzX3RyZW5kaW5nJmN0PWc/oYtVHSxngR3lC/giphy.gif",
	"https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExM3B2MmtvdW55YTRndnFyY2g5dGpxMzlmbmZucnR2czViOHphaGtpeCZlcD12MV9naWZzX3RyZW5kaW5nJmN0PWc/JpG2A9P3dPHXaTYrwu/giphy.gif",
	"https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExM3B2MmtvdW55YTRndnFyY2g5dGpxMzlmbmZucnR2czViOHphaGtpeCZlcD12MV9naWZzX3RyZW5kaW5nJmN0PWc/JpG2A9P3dPHXaTYrwu/giphy.gif",
	"https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExM3B2MmtvdW55YTRndnFyY2g5dGpxMzlmbmZucnR2czViOHphaGtpeCZlcD12MV9naWZzX3RyZW5kaW5nJmN0PWc/OuQmhmAAdJFLi/giphy.gif",
	"https://media.giphy.com/media/v1.Y2lkPWVjZjA1ZTQ3MjJzMHQ5d2hxN3JvZDhlem54NDYyeDJ0dzIxcm9qd2F2OGxoZ2loMyZlcD12MV9naWZzX3RyZW5kaW5nJmN0PWc/YaP3iYxN3T8nIEN5rD/giphy.gif",
	"https://media.giphy.com/media/v1.Y2lkPWVjZjA1ZTQ3MjJzMHQ5d2hxN3JvZDhlem54NDYyeDJ0dzIxcm9qd2F2OGxoZ2loMyZlcD12MV9naWZzX3RyZW5kaW5nJmN0PWc/dQU7k68XRCnnsUSoZO/giphy.gif",
}

// sampleTexts — pool of random texts for comments.
var sampleTexts = []string{
	"Wow, this is actually awesome!",
	"Don't agree, but interesting take.",
	"Hahaha, what a gif 😂",
	"Thanks for sharing!",
	"Where did you get this?",
	"Great video, not my first time watching.",
	"Something is clearly off here...",
	"Subscribed, top-tier content.",
	"First!",
	"Who else is watching in 2026?",
	"That moment at 3:21 is just 🔥",
	"I don't get the hype around this.",
	"Saved it, thanks.",
	"Damn, so well put.",
	"This is the best thing I've seen today.",
	"",
}

func RunCommentsFixture(userId *string, count *int, videoId *string, nestedChance *float64) {
	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.NewDatabase(c)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	ctx := context.Background()

	createdIds := make([]string, 0, *count)

	insertQuery := `
		INSERT INTO comments (user_id, video_id, gif_url, text, answer_to)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	for i := 0; i < *count; i++ {
		var gifUrl *string
		if len(gifUrls) > 0 {
			g := gifUrls[rand.Intn(len(gifUrls))]
			gifUrl = &g
		}

		text := sampleTexts[rand.Intn(len(sampleTexts))]

		// Make sure the comment has at least something: if there is no text and no gif, set a text.
		if text == "" && gifUrl == nil {
			text = "No comment 🙂"
		}

		var answerTo *string
		// A comment can be nested only if there is already something to reply to.
		if len(createdIds) > 0 && rand.Float64() < *nestedChance {
			parent := createdIds[rand.Intn(len(createdIds))]
			answerTo = &parent
		}

		var id string
		err := pool.QueryRow(ctx, insertQuery, *userId, *videoId, gifUrl, text, answerTo).Scan(&id)
		if err != nil {
			log.Fatalf("failed to create comment #%d: %v", i+1, err)
		}

		createdIds = append(createdIds, id)
	}

	log.Printf("created %d comments for user_id=%s, video_id=%s", len(createdIds), *userId, *videoId)
}
