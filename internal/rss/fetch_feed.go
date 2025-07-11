package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}

	headers := map[string]string{
		"User-Agent": "gator",
	}

	data, err := handleRequest(ctx, feedURL, client, headers)
	if err != nil {
		return nil, fmt.Errorf("http error: %v", err)
	}

	feed := RSSFeed{}

	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("error unmarshaling data: %v", err)
	}

	unescapeFeed(&feed)

	return &feed, nil
}

func handleRequest(ctx context.Context, feedURL string, client *http.Client, headers ...map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}

	for _, headerMap := range headers {
		for key, value := range headerMap {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return data, nil
}

func unescapeFeed(feed *RSSFeed) {
	feed.Channel.Title, feed.Channel.Description =
		html.UnescapeString(feed.Channel.Title), html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title, feed.Channel.Item[i].Description =
			html.UnescapeString(feed.Channel.Item[i].Title), html.UnescapeString(feed.Channel.Item[i].Description)
	}
}
