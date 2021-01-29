package telegram

import (
	"encoding/json"
	"fmt"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
	xurls "mvdan.cc/xurls/v2"
)

const (
	HackerNewsChatID = -1001020923877
)

type LinkSaver func([]string) error

func NewWebhookHandler(token, webhookUrl string, saveLinks LinkSaver) (*tb.Webhook, error) {
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

	b.Handle(tb.OnText, func(m *tb.Message) {
		links := getLinks(m)
		log.Printf("Links found: %v", links)

		if m.OriginalChat != nil && m.OriginalChat.ID == HackerNewsChatID {
			links = getHNCommentLinks(links)
			log.Printf("Message from HN channel, links filtered to: %v", links)
		}

		err := saveLinks(links)
		if err != nil {
			b.Reply(m, fmt.Sprintf("Error saving link: %s", err), tb.NoPreview)
		}

		b.Reply(m, "Saved.")
	})

	go b.Start()

	return webhook, nil
}

func getLinks(m *tb.Message) []string {
	return xurls.Relaxed().FindAllString(m.Text, -1)
}

const hnCommentLinkPrefix = "https://readhacker.news/c/"

func getHNCommentLinks(links []string) []string {
	for _, link := range links {
		if len(link) < len(hnCommentLinkPrefix) {
			continue
		}

		if link[:len(hnCommentLinkPrefix)] == hnCommentLinkPrefix {
			return []string{link}
		}
	}

	return []string{}
}
