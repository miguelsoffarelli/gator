package commands

import (
	"context"
	"fmt"
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
