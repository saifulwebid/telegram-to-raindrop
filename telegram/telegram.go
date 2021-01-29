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

type Settings struct {
	Token      string
	WebhookURL string
	LinkSaver  LinkSaver

	AdminUser int
}

func NewWebhookHandler(s Settings) (*tb.Webhook, error) {
	webhook := &tb.Webhook{
		Endpoint: &tb.WebhookEndpoint{
			PublicURL: s.WebhookURL,
		},
	}

	b, err := tb.NewBot(tb.Settings{
		Token:   s.Token,
		Verbose: true,
	})
	if err != nil {
		return nil, fmt.Errorf("telebot.NewBot failed: %w", err)
	}

	b.Poller = webhook

	if s.AdminUser != 0 {
		b.Poller = tb.NewMiddlewarePoller(b.Poller, func(u *tb.Update) bool {
			if u.Message == nil {
				return true
			}

			if u.Message.Sender == nil {
				return true
			}

			if u.Message.Sender.ID != s.AdminUser {
				admin := &tb.User{ID: s.AdminUser}
				b.Send(admin, fmt.Sprintf(
					"Unknown user: %s %s (username: %s)",
					u.Message.Sender.FirstName,
					u.Message.Sender.LastName,
					u.Message.Sender.Username,
				))
				b.Forward(admin, u.Message)

				return false
			}

			return true
		})
	}

	b.Poller = tb.NewMiddlewarePoller(b.Poller, func(u *tb.Update) bool {
		j, _ := json.MarshalIndent(u, "", "\t")
		log.Printf("Update received:\n%v", string(j))
		return true
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		links := getLinks(m)
		log.Printf("Links found: %v", links)

		if m.OriginalChat != nil && m.OriginalChat.ID == HackerNewsChatID {
			links = getHNCommentLinks(links)
			log.Printf("Message from HN channel, links filtered to: %v", links)
		}

		if len(links) == 0 {
			b.Reply(m, "No links detected.")
			return
		}

		err := s.LinkSaver(links)
		if err != nil {
			b.Reply(m, fmt.Sprintf("Error saving link: %s", err), tb.NoPreview)
			return
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
