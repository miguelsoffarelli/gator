package commands

import (
	"context"
	"fmt"

	"github.com/miguelsoffarelli/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		ctx := context.Background()
		user, err := s.Db.GetUserByName(ctx, s.Cfg.Current_user_name)
		if err != nil {
			return fmt.Errorf("error retrieving user from database: %v", err)
		}

		return handler(s, cmd, user)
	}
}
