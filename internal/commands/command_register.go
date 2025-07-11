package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/miguelsoffarelli/go-blog-aggregator/internal/database"
)

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: username required")
	}

	name := cmd.Args[0]
	ctx := context.Background()

	_, err := s.Db.GetUser(ctx, name)
	if err == nil {
		return fmt.Errorf("error: user already exists")
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("database error: %v", err)
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	user, err := s.Db.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("error creating user %v: %v", name, err)
	}

	s.Cfg.SetUser(user.Name)

	fmt.Printf("User %v successfully created!\n", user.Name)
	log.Printf("ID: %v\n", user.ID)
	log.Printf("Created At: %v\n", user.CreatedAt)
	log.Printf("Updated At: %v\n", user.UpdatedAt)
	log.Printf("Name: %v\n", user.Name)

	return nil
}
