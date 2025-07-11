package commands

import (
	"context"
	"database/sql"
	"fmt"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: username required")
	}

	name := cmd.Args[0]

	if _, err := s.Db.GetUserByName(context.Background(), name); err == sql.ErrNoRows {
		return fmt.Errorf("error: user doesn't exist")
	} else if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	if err := s.Cfg.SetUser(name); err != nil {
		return err
	}

	fmt.Printf("User %v has been successfully set\n", cmd.Args[0])

	return nil
}
