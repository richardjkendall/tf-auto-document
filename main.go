package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardjkendall/tf-auto-document/parser"
	"github.com/richardjkendall/tf-auto-document/scangit"
)

func scanModulesFolder(path string, scanner *scangit.ScanGit) ([]parser.ModuleDetails, error) {
	var r []parser.ModuleDetails
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return r, err
	}
	for _, file := range files {
		// ignore . files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		fullPath := filepath.Join(path, file.Name())
		// only look at directories
		if file.IsDir() {
			m, err := parser.New().ParseModule(fullPath)
			if err != nil {
				return r, err
			}
			r = append(r, m)
		}
	}
	return r, nil
}

func main() {

	folderToScan := os.Args[1]
	fmt.Println("Working on repository: " + folderToScan)

	// create gitscanner for this repo
	scanner := scangit.New()
	err := scanner.Open(folderToScan)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// load tags
	err = scanner.LoadTags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	/*
		c, err := scanner.GetCommits("modules/ecs-service/")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("commits %+v\n", c)*/

	mod, err := scanModulesFolder(folderToScan+"/modules", scanner)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("module details which came back: %+v", mod)

}
