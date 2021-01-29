package ttr

import (
	"log"
	"net/http"
	"os"
	"strconv"

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

	adminUser, err := strconv.Atoi(os.Getenv("TELEGRAM_ADMIN_ID"))
	if err != nil {
		log.Println("TELEGRAM_ADMIN_ID is empty or not a valid integer; disabling the feature")
		adminUser = 0
	} else {
		log.Printf("Admin user ID: %d", adminUser)
	}

	webhook, err = telegram.NewWebhookHandler(telegram.Settings{
		Token:      os.Getenv("TELEGRAM_TOKEN"),
		WebhookURL: os.Getenv("TELEGRAM_HOOK_URL"),
		LinkSaver:  raindropClient.Save,
		AdminUser:  adminUser,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func CFHandler(w http.ResponseWriter, r *http.Request) {
	webhook.ServeHTTP(w, r)
}
