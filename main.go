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

func createModuleReadme(path string, details CombinedModuleDetails) error {
	w := writer.New(path + "/README.md")
	w.H1Underline(details.TFDetails.Title)
	w.P(details.TFDetails.Desc)
	if len(details.TFDetails.Depends) > 0 {
		w.H2Underline("Depends on")
		for _, d := range details.TFDetails.Depends {
			w.Bullet(writer.MakeLink(d, "../"+d+"/README.md"))
		}
		w.P("")
	}
	if len(details.TFDetails.Partners) > 0 {
		w.H2Underline("Works with")
		for _, d := range details.TFDetails.Partners {
			w.Bullet(writer.MakeLink(d, "../"+d+"/README.md"))
		}
		w.P("")
	}
	w.H2Underline("Releases")
	var commitRows [][]string
	for _, commit := range details.GitDetails {
		if commit.Tag != "" {
			row := []string{commit.Tag, strings.Trim(commit.Message, "\r\n"), writer.InlineCode(commit.Hash[0:6])}
			commitRows = append(commitRows, row)
		}
	}
	if len(commitRows) > 0 {
		commitHeaders := []string{"Tag", "Message", "Commit"}
		w.Table(commitHeaders, commitRows)
	} else {
		w.P("There have been no releases yet for this module")
	}
	var varRows [][]string
	for _, variable := range details.TFDetails.Variables {
		dt := "`not specified`"
		if variable.DataType != "" {
			dt = writer.InlineCode(variable.DataType)
		}
		row := []string{writer.InlineCode(variable.Name), dt, variable.Desc, writer.InlineCode(variable.Def)}
		varRows = append(varRows, row)
	}
	varHeaders := []string{"Name", "Type", "Description", "Default Value"}
	w.H2Underline("Variables")
	w.Table(varHeaders, varRows)
	return w.WriteFile()
}

func createRootReadme(path string, details []CombinedModuleDetails) error {
	w := writer.New(path + "/README.md")
	w.H1Underline("Terraform Modules")
	w.P("This is a collection of terraform modules")
	w.P("Click on the links to see the details of each of the modules")
	w.P("This documentation is auto-generated from the terraform files using tf-auto-document.")
	var modRows [][]string
	for _, module := range details {
		if module.TFDetails.Title != "" {
			row := []string{module.TFDetails.Title, module.TFDetails.Desc, writer.MakeLink("more details", module.Folder+"/README.md")}
			modRows = append(modRows, row)
		} else {
			fmt.Printf("error no details found for module in %s\n", module.Folder)
		}
	}
	headers := []string{"Module", "Description", "Link"}
	w.H2Underline("Modules")
	w.Table(headers, modRows)
	return w.WriteFile()
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
			cmd.Folder = modulesfolder + "/" + file.Name()
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

	// create root md file
	rerr := createRootReadme(folderToScan, mod)
	if rerr != nil {
		fmt.Println(rerr)
		os.Exit(1)
	}

	// create each module's md file
	for _, m := range mod {
		merr := createModuleReadme(folderToScan+"/"+m.Folder, m)
		if merr != nil {
			fmt.Println(rerr)
			os.Exit(1)
		}
	}
}
