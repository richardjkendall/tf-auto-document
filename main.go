package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	LogLevel string `hcl:"log_level"`
}

func main() {
	var config Config

	folderToScan := os.Args[1]
	fmt.Println("Working on folder: " + folderToScan)

	err := hclsimple.DecodeFile(folderToScan, nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	log.Printf("Configuration is %#v", config)
}
