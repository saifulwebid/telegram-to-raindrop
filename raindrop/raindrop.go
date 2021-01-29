package raindrop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type Client struct {
	client *http.Client
}

func NewClient(token string) *Client {
	ctx := context.Background()

	client := &Client{
		client: oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: token,
		})),
	}

	return client
}

type Raindrop struct {
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Link    string `json:"link"`
}

func (c *Client) ParseLink(link string) (Raindrop, error) {
	params := url.Values{}
	params.Add("url", link)

	parseResponse, err := c.client.Get("https://api.raindrop.io/rest/v1/import/url/parse?" + params.Encode())
	if err != nil {
		return Raindrop{}, fmt.Errorf("failed to request URL parsing for %s: %w", link, err)
	}
	defer parseResponse.Body.Close()

	var parsedResult struct {
		Result       bool   `json:"result"`
		ErrorMessage string `json:"errorMessage"`
		Item         struct {
			Raindrop
			Meta struct {
				Canonical string `json:"canonical"`
			} `json:"meta"`
		} `json:"item"`
	}
	if err := json.NewDecoder(parseResponse.Body).Decode(&parsedResult); err != nil {
		log.Printf("response code: %v", parseResponse.Status)
		return Raindrop{}, fmt.Errorf("failed to decode parsed URL %s: %w", link, err)
	}
	if parsedResult.ErrorMessage != "" {
		return Raindrop{}, fmt.Errorf("failed to parse URL %s: %s", link, parsedResult.ErrorMessage)
	}

	rd := parsedResult.Item.Raindrop
	if parsedResult.Item.Meta.Canonical != "" {
		rd.Link = parsedResult.Item.Meta.Canonical
	} else {
		rd.Link = link
	}

	return rd, nil
}

func (c *Client) SaveRaindrop(rd Raindrop) error {
	rdToPost := struct {
		Raindrop
		PleaseParse struct{} `json:"pleaseParse"`
	}{
		Raindrop:    rd,
		PleaseParse: struct{}{},
	}

	raindropJson, err := json.Marshal(rdToPost)
	if err != nil {
		return fmt.Errorf("failed to encode JSON to be posted: %w", err)
	}

	postResponse, err := c.client.Post("https://api.raindrop.io/rest/v1/raindrop", "application/json", bytes.NewReader(raindropJson))
	if err != nil {
		return fmt.Errorf("failed to post URL %s: %w", rd.Link, err)
	}
	defer postResponse.Body.Close()

	var respMap map[string]interface{}
	if err := json.NewDecoder(postResponse.Body).Decode(&respMap); err != nil {
		return fmt.Errorf("failed to parse response body fom Raindrop.io on posting %s: %w", rd.Link, err)
	}
	if !respMap["result"].(bool) {
		return fmt.Errorf("failed to post link %s: %s", rd.Link, respMap["errorMessage"].(string))
	}

	return nil
}

func (c *Client) Save(links []string) error {
	for _, link := range links {
		rd, err := c.ParseLink(link)
		if err != nil {
			return err
		}

		err = c.SaveRaindrop(rd)
		if err != nil {
			return err
		}
	}

	return nil
}
