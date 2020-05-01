package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// Parser is the object which holds the methods needed to scan hcl files looking for the details we need
type Parser struct {
	hclParser *hclparse.Parser
}

// ModuleDetails contains the details of the module being scanned
type ModuleDetails struct {
	title    string
	desc     string
	partners []string
	depends  []string
}

// New creates a new instance of Parser
func New() *Parser {
	return &Parser{
		hclParser: hclparse.NewParser(),
	}
}

// ParseFolder opens a folder and finds the *.tf files so we can scan them
// and get the information needed to write the documentation files
func (parser *Parser) ParseFolder(path string) (*hcl.Blocks, *hcl.Blocks, *hcl.Blocks, error) {
	// parse all the terraform files I find in the directory
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, nil, nil, err
	}
	for _, file := range files {
		// ignore . files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		fullPath := filepath.Join(path, file.Name())
		// ignore directories
		if !file.IsDir() {
			if strings.HasSuffix(file.Name(), ".tf") {
				fmt.Println("\tFile: " + fullPath)
				if strings.HasSuffix(file.Name(), "main.tf") {
					details, err := parser.getMainDetails(fullPath)
					if err != nil {
						return nil, nil, nil, err
					}
					fmt.Printf("%+v\n", details)
				}
				_, diagnostics := parser.hclParser.ParseHCLFile(fullPath)
				if diagnostics != nil && diagnostics.HasErrors() {
					return nil, nil, nil, diagnostics
				}
			}
		}
	}

	// run parser on all files
	var blocks hcl.Blocks
	for _, file := range parser.hclParser.Files() {
		fileBlocks, err := parser.parseFile(file)
		if err != nil {
			return nil, nil, nil, err
		}
		blocks = append(blocks, fileBlocks...)
	}

	for _, block := range blocks {
		fmt.Println(block.Type)
	}

	return nil, nil, nil, nil
}

func (parser *Parser) parseFile(file *hcl.File) (hcl.Blocks, error) {
	contents, diagnostics := file.Body.Content(terraformSchema)
	if diagnostics != nil && diagnostics.HasErrors() {
		return nil, diagnostics
	}
	if contents == nil {
		return nil, fmt.Errorf("File is empty")
	}
	return contents.Blocks, nil
}

// getMainDetails scans a tf file looking for a specific pattern of comment which contains the details of the file
// outputs a struct containing these details
func (parser *Parser) getMainDetails(path string) (ModuleDetails, error) {
	var r ModuleDetails
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return r, err
	}
	var re = regexp.MustCompile(`(?m)^\/\*\ntitle:\s+([\w\-]+)\ndesc:\s+([\w\-\s\.,]+)\n(partners:\s+[\w\-,\s]+\n)?(depends:\s+[\w\-,\s]+\n)?\*\/`)
	match := re.FindAllStringSubmatch(string(data), -1)
	var title string
	var desc string
	var partners []string
	var depends []string
	if match == nil {
		return r, nil
	}
	title = match[0][1]
	desc = match[0][2]
	for i := 3; i < len(match[0]); i++ {
		temp := match[0][i]
		if strings.HasPrefix(temp, "partners:") {
			partners = trimAll(strings.Split(strings.Split(temp, ":")[1], ","))
		}
		if strings.HasPrefix(temp, "depends:") {
			depends = trimAll(strings.Split(strings.Split(temp, ":")[1], ","))
		}
	}
	r = ModuleDetails{
		title:    title,
		desc:     desc,
		partners: partners,
		depends:  depends,
	}
	return r, nil
}

// trimAll takes the elements of a slice of strings and trims all the whitespace off the strings in the slice
func trimAll(input []string) []string {
	output := make([]string, len(input))
	for i, s := range input {
		output[i] = strings.Trim(s, " ")
	}
	return output
}
