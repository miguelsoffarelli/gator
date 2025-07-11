package commands

import (
	"context"
	"fmt"
)

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users from database: %v", err)
	}

	for _, user := range users {
		if user.Name == s.Cfg.Current_user_name {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}

	return nil
}
