package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/miguelsoffarelli/gator/internal/database"
)

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("must provide url")
	}

	ctx := context.Background()
	url := cmd.Args[0]

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

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 0 {
		fmt.Println("Parameters ignored: command following takes no parameters")
	}

	ctx := context.Background()

	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting feeds from database: %v", err)
	}

	for _, feed := range feedFollows {
		fmt.Println(feed.FeedName)
	}
	fmt.Println("=====================================================")

	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("must provide url from feed to unfollow")
	}

	url := cmd.Args[0]
	ctx := context.Background()
	params := database.UnfollowParams{
		UserID: user.ID,
		Url:    url,
	}

	if err := s.Db.Unfollow(ctx, params); err != nil {
		return fmt.Errorf("error unfollowing feed: %v", err)
	}

	return nil
}
