package main

import (
	"log"
	"os"

	config "github.com/miguelsoffarelli/go-blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	s := state{}
	s.cfg = &cfg

	c := commands{
		cmds: make(map[string]func(*state, command) error),
	}

	c.register("login", handlerLogin)

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("Error: not enough arguments")
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	if err := c.run(&s, cmd); err != nil {
		log.Fatalf("%v", err)
	}
}
