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

func scanModulesFolder(path string, modulesfolder string, scanner *scangit.ScanGit) ([]parser.ModuleDetails, error) {
	var r []parser.ModuleDetails
	files, err := ioutil.ReadDir(path + "/" + modulesfolder)
	if err != nil {
		return r, err
	}
	for _, file := range files {
		// ignore . files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		fmt.Printf("folder = %s\n", modulesfolder+"/"+file.Name())
		fullPath := filepath.Join(path+"/"+modulesfolder, file.Name())
		// only look at directories
		if file.IsDir() {
			m, err := parser.New().ParseModule(fullPath)
			if err != nil {
				return r, err
			}
			r = append(r, m)
			// need to get commits
			c, err := scanner.GetCommits(modulesfolder + "/" + file.Name())
			if err != nil {
				return r, err
			}
			fmt.Printf("commits %+v\n", c)
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

	// scan terraform files
	mod, err := scanModulesFolder(folderToScan, "modules", scanner)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("module details which came back: %+v", mod)

}
