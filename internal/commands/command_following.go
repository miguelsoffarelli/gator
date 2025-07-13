package commands

import (
	"context"
	"fmt"
)

func HandlerFollowing(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		fmt.Println("Parameters ignored: command following takes no parameters")
	}

	ctx := context.Background()

	currentUser, err := s.Db.GetUserByName(ctx, s.Cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	currentUserID := currentUser.ID

	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, currentUserID)
	if err != nil {
		return fmt.Errorf("error getting feeds from database: %v", err)
	}

	for _, feed := range feedFollows {
		fmt.Println(feed.FeedName)
	}
	fmt.Println("=====================================================")

	return nil
}
