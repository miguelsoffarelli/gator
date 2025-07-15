package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	commands "github.com/miguelsoffarelli/gator/internal/commands"
	config "github.com/miguelsoffarelli/gator/internal/config"
	database "github.com/miguelsoffarelli/gator/internal/database"
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

	s := commands.State{}
	s.Cfg, s.Db = &cfg, dbQueries

	c := commands.Commands{
		Cmds: make(map[string]func(*commands.State, commands.Command) error),
	}

	c.Register("login", commands.HandlerLogin)
	c.Register("register", commands.HandlerRegister)
	c.Register("reset", commands.HandlerReset)
	c.Register("users", commands.HandlerUsers)
	c.Register("agg", commands.HandlerAgg)
	c.Register("addfeed", commands.MiddlewareLoggedIn(commands.HandlerAddFeed))
	c.Register("feeds", commands.HandlerFeeds)
	c.Register("follow", commands.MiddlewareLoggedIn(commands.HandlerFollow))
	c.Register("following", commands.MiddlewareLoggedIn(commands.HandlerFollowing))
	c.Register("unfollow", commands.MiddlewareLoggedIn(commands.HandlerUnfollow))
	c.Register("browse", commands.MiddlewareLoggedIn(commands.HandlerBrowse))

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("Error: not enough arguments")
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	cmd := commands.Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	if err := c.Run(&s, cmd); err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}
}
