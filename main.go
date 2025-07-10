package main

import (
	"encoding/json"
	"fmt"

	config "github.com/miguelsoffarelli/go-blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config file: %v", err)
	}

	cfg.SetUser("miguel")

	cfg, err = config.Read()
	if err != nil {
		fmt.Printf("Error reading config file: %v", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling json: %v", err)
	}

	fmt.Println(string(data))
}
