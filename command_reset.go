package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		return fmt.Errorf("error resetting database: %v", err)
	}

	fmt.Println("Database successfully reset.")
	return nil
}
