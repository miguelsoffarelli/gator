package commands

import (
	"context"
	"fmt"

	"github.com/miguelsoffarelli/go-blog-aggregator/internal/rss"
)

func HandlerAgg(s *State, cmd Command) error {
	ctx := context.Background()
	URL := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(ctx, URL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil
}
