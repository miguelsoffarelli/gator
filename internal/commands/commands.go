package commands

import (
	"fmt"

	config "github.com/miguelsoffarelli/gator/internal/config"
	database "github.com/miguelsoffarelli/gator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Cmds map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if _, ok := c.Cmds[cmd.Name]; !ok {
		return fmt.Errorf("command %v not found", cmd.Name)
	}

	if err := c.Cmds[cmd.Name](s, cmd); err != nil {
		return err
	}

	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Cmds[name] = f
}
