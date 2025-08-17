package main

import (
	"log/slog"
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
		filename, err := bot.GenerateImage()
		if err != nil {
			slog.Error("Error generating image:", "ERROR", err)
			return
		}
		fileKey, err := bot.UploadImage(baseUrl, channelID, token, filename)
		if err != nil {
			slog.Error("Error uploading image:", "ERROR", err)
			return
		}
		err = bot.PostMessage(baseUrl, channelID, token, fileKey)
		if err != nil {
			slog.Error("Error posting message:", "ERROR", err)
			return
		}
	}

	c := cron.New()
	c.AddFunc("0 6 * * *", job)
	c.Start()

	select {}
}
