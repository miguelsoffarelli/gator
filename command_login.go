package main

import (
	"context"
	"database/sql"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("error: username required")
	}

	name := cmd.args[0]

	if _, err := s.db.GetUser(context.Background(), name); err == sql.ErrNoRows {
		return fmt.Errorf("error: user doesn't exist")
	} else if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	if err := s.cfg.SetUser(name); err != nil {
		return err
	}

	fmt.Printf("User %v has been successfully set\n", cmd.args[0])

	return nil
}
