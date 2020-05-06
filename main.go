package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardjkendall/tf-auto-document/parser"
	"github.com/richardjkendall/tf-auto-document/scangit"
	"github.com/richardjkendall/tf-auto-document/writer"
)

// CombinedModuleDetails holds the combined module details
type CombinedModuleDetails struct {
	Folder     string
	TFDetails  parser.ModuleDetails
	GitDetails []scangit.GitCommit
}

func createRootReadme(path string, details []CombinedModuleDetails) error {
	w := writer.New(path + "/README.md")
	w.H1Underline("Terraform Modules")
	w.P("This is a collection of terraform modules")
	w.P("Click on the links to see the details of each of the modules")
	for _, module := range details {

	}
}

func scanModulesFolder(path string, modulesfolder string, scanner *scangit.ScanGit) ([]CombinedModuleDetails, error) {
	var r []CombinedModuleDetails
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
			var cmd CombinedModuleDetails
			m, err := parser.New().ParseModule(fullPath)
			if err != nil {
				return r, err
			}
			cmd.TFDetails = m

			// need to get commits
			c, err := scanner.GetCommits(modulesfolder + "/" + file.Name())
			if err != nil {
				return r, err
			}
			cmd.GitDetails = c
			r = append(r, cmd)
			//fmt.Printf("commits %+v\n", c)
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
