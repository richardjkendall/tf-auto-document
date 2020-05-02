package main

import (
	"fmt"
	"os"

	"github.com/richardjkendall/tf-auto-document/parser"
)

func main() {

	folderToScan := os.Args[1]
	fmt.Println("Working on folder: " + folderToScan)

	_, err := parser.New().ParseFolder(folderToScan)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
