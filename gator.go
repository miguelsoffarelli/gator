package main

import (
	"fmt"

	config "github.com/miguelsoffarelli/go-blog-aggregator/internal/config"
	database "github.com/miguelsoffarelli/go-blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if _, ok := c.cmds[cmd.name]; !ok {
		return fmt.Errorf("command %v not found", cmd.name)
	}

	if err := c.cmds[cmd.name](s, cmd); err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
