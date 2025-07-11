package commands

import (
	"context"
	"fmt"
)

func HandlerReset(s *State, cmd Command) error {
	if err := s.Db.Reset(context.Background()); err != nil {
		return fmt.Errorf("error resetting database: %v", err)
	}

	fmt.Println("Database successfully reset.")
	return nil
}
