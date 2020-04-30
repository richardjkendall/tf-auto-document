package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Parser struct {
	hclParser *hclparse.Parser
}

func New() *Parser {
	return &Parser{
		hclParser: hclparse.NewParser(),
	}
}

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
