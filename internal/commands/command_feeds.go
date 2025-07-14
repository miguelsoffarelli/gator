package commands

import (
	"context"
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
	ctx := context.Background()
	URL := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(ctx, URL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil
}
