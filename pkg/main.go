package main

import (
	"os"

	"github.com/claustra01/typetalk-progress-bar-bot/pkg/bot"
	"github.com/robfig/cron/v3"
)

func main() {
	topicId := os.Getenv("TYPETALK_TOPIC_ID")
	if topicId == "" {
		panic("TYPETALK_TOPIC_ID is not set")
	}

	token := os.Getenv("TYPETALK_TOKEN")
	if token == "" {
		panic("TYPETALK_TOKEN is not set")
	}

	job := func() {
		filename := bot.GenerateImage()
		fileKey := bot.UploadImage(topicId, token, filename)
		bot.PostMessage(topicId, token, fileKey)
	}

	c := cron.New()
	c.AddFunc("0 6 * * *", job)
	c.Start()

	select {}
}
