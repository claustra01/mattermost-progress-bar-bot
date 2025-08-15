package main

import (
	"os"

	"github.com/claustra01/mattermost-progress-bar-bot/pkg/bot"
	"github.com/robfig/cron/v3"
)

func main() {
	baseUrl := os.Getenv("MATTERMOST_BASE_URL")
	if baseUrl == "" {
		panic("MATTERMOST_BASE_URL is not set")
	}

	channelID := os.Getenv("MATTERMOST_CHANNEL_ID")
	if channelID == "" {
		panic("MATTERMOST_CHANNEL_ID is not set")
	}

	token := os.Getenv("MATTERMOST_TOKEN")
	if token == "" {
		panic("MATTERMOST_TOKEN is not set")
	}

	job := func() {
		filename := bot.GenerateImage()
		fileKey := bot.UploadImage(baseUrl, channelID, token, filename)
		bot.PostMessage(baseUrl, channelID, token, fileKey)
	}

	c := cron.New()
	c.AddFunc("0 6 * * *", job)
	c.Start()

	select {}
}
