package telegram

import (
	"encoding/json"
	"fmt"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

func NewWebhookHandler(token, webhookUrl string) (*tb.Webhook, error) {
	webhook := &tb.Webhook{
		Endpoint: &tb.WebhookEndpoint{
			PublicURL: webhookUrl,
		},
	}

	loggedWebhook := tb.NewMiddlewarePoller(webhook, func(u *tb.Update) bool {
		j, _ := json.MarshalIndent(u, "", "\t")
		log.Printf("Update received:\n%v", string(j))
		return true
	})

	b, err := tb.NewBot(tb.Settings{
		Token:   token,
		Poller:  loggedWebhook,
		Verbose: true,
	})
	if err != nil {
		return nil, fmt.Errorf("telebot.NewBot failed: %w", err)
	}

	go b.Start()

	return webhook, nil
}
