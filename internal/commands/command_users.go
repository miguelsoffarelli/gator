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

func HandlerUsers(s *State, cmd Command) error {
	if len(cmd.Args) != 0 {
		fmt.Println("Parameters ignored: command users takes no parameters")
	}

	ctx := context.Background()

	users, err := s.Db.ListUsers(ctx)
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

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: username required")
	}

	name := cmd.Args[0]
	ctx := context.Background()

	_, err := s.Db.GetUserByName(ctx, name)
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
