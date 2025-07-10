package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	config "github.com/miguelsoffarelli/go-blog-aggregator/internal/config"
	database "github.com/miguelsoffarelli/go-blog-aggregator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	dbQueries := database.New(db)

	s := state{}
	s.cfg, s.db = &cfg, dbQueries

	c := commands{
		cmds: make(map[string]func(*state, command) error),
	}

	c.register("login", handlerLogin)
	c.register("register", handlerRegister)

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
		os.Exit(1)
	}
}
