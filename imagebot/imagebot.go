package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	bucket := "line-bot-01-269403.appspot.com"

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					object := message.Text + ".jpg"

					ctx, cancel := context.WithTimeout(ctx, time.Second*10)
					defer cancel()
					obj := client.Bucket(bucket).Object(object)
					if _, err := obj.Attrs(ctx); err != nil {
						if _, err = bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage("錯誤，請重新上傳")).Do(); err != nil {
							log.Fatal(err)
						}
						log.Fatalf("Message: %v", err)
						return
					}

					originalContentURL := "https://storage.googleapis.com/line-bot-01-269403.appspot.com/" + message.Text + ".jpg"
					previewImageURL := "https://storage.googleapis.com/line-bot-01-269403.appspot.com/" + message.Text + ".jpg"

					if _, err = bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewImageMessage(originalContentURL, previewImageURL)).Do(); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
