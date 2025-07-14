package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miguelsoffarelli/go-blog-aggregator/internal/database"
	"github.com/miguelsoffarelli/go-blog-aggregator/internal/rss"
)

func HandlerFeeds(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		fmt.Println("Parameters ignored: command feeds takes no parameters")
	}

	ctx := context.Background()
	feeds, err := s.Db.ListFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error getting feeds from database: %v", err)
	}

	for _, feed := range feeds {
		name, URL := feed.Name, feed.Url
		user, err := s.Db.GetUserById(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting name of the user that created the feed: %v", err)
		}
		userName := user.Name

		fmt.Printf("Name: %v\n", name)
		fmt.Printf("URL: %v\n", URL)
		fmt.Printf("Created By: %v\n", userName)
		fmt.Println("==============================================")
	}

	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("must provide feed name and url")
	}

	feedName, feedURL := cmd.Args[0], cmd.Args[1]

	ctx := context.Background()

	createFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	feed, err := s.Db.CreateFeed(ctx, createFeedParams)
	if err != nil {
		return fmt.Errorf("error adding feed to database: %v", err)
	}

	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	if _, err = s.Db.CreateFeedFollow(ctx, createFeedFollowParams); err != nil {
		return fmt.Errorf("error following created feed: %v", err)
	}

	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("Created At: %v\n", feed.CreatedAt)
	fmt.Printf("Updated At: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("URL: %v\n", feed.Url)
	fmt.Printf("User ID: %v\n", feed.UserID)

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("must provide time between requests argument")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing time between requests: %v", err)
	}

	fmt.Printf("Collecting feeds every %v", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("error scraping feed: %v", err)
		}
	}
}

func scrapeFeeds(s *State) error {
	ctx := context.Background()
	nextFeed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %v", err)
	}

	params := database.MarkFeedFetchedParams{
		ID: nextFeed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	if err := s.Db.MarkFeedFetched(ctx, params); err != nil {
		return fmt.Errorf("error marking feed as fetched: %v", err)
	}

	feed, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	for _, item := range feed.Channel.Item {
		fmt.Println(item.Title)
	}

	return nil
}
