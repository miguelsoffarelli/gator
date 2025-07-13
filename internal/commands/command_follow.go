package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miguelsoffarelli/go-blog-aggregator/internal/database"
)

func HandlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("must provide url")
	}

	ctx := context.Background()
	url := cmd.Args[0]

	user, err := s.Db.GetUserByName(ctx, s.Cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("error getting current user from the database: %v", err)
	}

	feed, err := s.Db.GetFeedByUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("error getting feed from database: %v", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	followedFeed, err := s.Db.CreateFeedFollow(ctx, params)
	if err != nil {
		return fmt.Errorf("error following feed: %v", err)
	}

	fmt.Println("Feed successfully followed!")
	fmt.Println("")
	fmt.Printf("Feed Name: %v\n", followedFeed.FeedName)
	fmt.Printf("Current User: %v\n", followedFeed.UserName)
	fmt.Println("=============================================")

	return nil
}
