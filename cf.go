package ttr

import (
	"log"
	"net/http"
	"os"

	"github.com/saifulwebid/telegram-to-raindrop/raindrop"
	"github.com/saifulwebid/telegram-to-raindrop/telegram"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	webhook *tb.Webhook
)

func init() {
	var err error

	raindropClient := raindrop.NewClient(os.Getenv("RAINDROP_TOKEN"))

	webhook, err = telegram.NewWebhookHandler(os.Getenv("TELEGRAM_TOKEN"), os.Getenv("TELEGRAM_HOOK_URL"), raindropClient.Save)
	if err != nil {
		log.Fatalln(err)
	}
}

func CFHandler(w http.ResponseWriter, r *http.Request) {
	webhook.ServeHTTP(w, r)
}
