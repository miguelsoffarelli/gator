package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miguelsoffarelli/go-blog-aggregator/internal/database"
)

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("must provide feed name and url")
	}

	feedName, feedURL := cmd.Args[0], cmd.Args[1]

	ctx := context.Background()

	user, err := s.Db.GetUserByName(ctx, s.Cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("error getting current user: %v", err)
	}

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	feed, err := s.Db.CreateFeed(ctx, params)
	if err != nil {
		return fmt.Errorf("error adding feed to database: %v", err)
	}

	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("Created At: %v\n", feed.CreatedAt)
	fmt.Printf("Updated At: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("URL: %v\n", feed.Url)
	fmt.Printf("User ID: %v\n", feed.UserID)

	return nil
}
