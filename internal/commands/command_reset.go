package commands

import (
	"context"
	"fmt"
)

func HandlerReset(s *State, cmd Command) error {
	if err := s.Db.Reset(context.Background()); err != nil {
		return fmt.Errorf("error resetting database: %v", err)
	}

	if err := s.Cfg.SetUser(""); err != nil {
		return fmt.Errorf("error resetting current user name: %v", err)
	}

	fmt.Println("Database successfully reset.")
	return nil
}
