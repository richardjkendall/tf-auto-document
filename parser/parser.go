package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclparse"

	"github.com/zclconf/go-cty/cty"
)

// Parser is the object which holds the methods needed to scan hcl files looking for the details we need
type Parser struct {
	hclParser *hclparse.Parser
}

// ModuleDetails contains the details of the module being scanned
type ModuleDetails struct {
	Title     string
	Desc      string
	Partners  []string
	Depends   []string
	Variables []VariableDetails
	Outputs   []OutputDetails
}

// VariableDetails contains the details of the variables defined by the module
type VariableDetails struct {
	Name     string
	Desc     string
	Def      string
	DataType string
}

// OutputDetails contains the details of the outputs defined by the module
type OutputDetails struct {
	Name string
	Desc string
}

// New creates a new instance of Parser
func New() *Parser {
	return &Parser{
		hclParser: hclparse.NewParser(),
	}
}

// ParseFolder opens a folder and finds the *.tf files so we can scan them
// and get the information needed to write the documentation files
func (parser *Parser) ParseFolder(path string) (ModuleDetails, error) {
	var r ModuleDetails
	var v []VariableDetails
	var o []OutputDetails
	// parse all the terraform files I find in the directory
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
		// ignore directories
		if !file.IsDir() {
			// only look at Terraform files
			if strings.HasSuffix(file.Name(), ".tf") {
				fmt.Println("\tFile: " + fullPath)
				// check the main.tf file for the comments
				if strings.HasSuffix(file.Name(), "main.tf") {
					r, err := parser.getMainDetails(fullPath)
					if err != nil {
						return r, err
					}
				}
				_, diagnostics := parser.hclParser.ParseHCLFile(fullPath)
				if diagnostics != nil && diagnostics.HasErrors() {
					return r, diagnostics
				}
			}
		}
	}

	// run parser on all files
	var blocks hcl.Blocks
	for _, file := range parser.hclParser.Files() {
		fileBlocks, err := parser.parseFile(file)
		if err != nil {
			return r, err
		}
		blocks = append(blocks, fileBlocks...)
	}

	// go through the variables
	ctx := &hcl.EvalContext{}
	for _, block := range blocks.OfType("variable") {
		var varDetails VariableDetails
		variableName := block.Labels[0]
		varDetails.Name = variableName
		fmt.Printf("\t\tvariable name: %s\n", variableName)
		// go through the attributes of the variable
		attributes, diagnostics := block.Body.JustAttributes()
		if diagnostics != nil && diagnostics.HasErrors() {
			return r, diagnostics
		}
		for _, attribute := range attributes {
			val, _ := attribute.Expr.Value(ctx)
			// get data type
			if attribute.Name == "type" {
				valType, err := typeexpr.Type(attribute.Expr)
				if err != nil {
					return r, err
				}
				varDetails.DataType = typeexpr.TypeString(valType)
			}
			// get description
			if attribute.Name == "description" && val.Type() == cty.String {
				varDetails.Desc = val.AsString()
			}
			// get default
			if attribute.Name == "default" {
				varDetails.Def = convertValueToString(val)
			}

		}
		v = append(v, varDetails)
	}
	r.Variables = v

	// go through the outputs if they are present
	for _, block := range blocks.OfType("output") {
		var outDetails OutputDetails
		outputName := block.Labels[0]
		outDetails.Name = outputName
		// find the description attribute if it is present
		attributes, diagnostics := block.Body.JustAttributes()
		if diagnostics != nil && diagnostics.HasErrors() {
			return r, diagnostics
		}
		for _, attribute := range attributes {
			val, _ := attribute.Expr.Value(ctx)
			if attribute.Name == "description" {
				outDetails.Desc = val.AsString()
			}
		}
		o = append(o, outDetails)
	}
	r.Outputs = o

	//fmt.Printf("%+v\n", r)
	return r, nil
}

// parseFile gets the contents of the file for later use
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
	var re = regexp.MustCompile(`(?m)^\/\*\r?\ntitle:\s+([\w\-]+)\r?\ndesc:\s+([\w\-\t \.,\/<>="';!@#$%^&*()_+~:]+)\r?\n(partners:\s+[\w\-,\s]+\r?\n)?(depends:\s+[\w\-,\s]+\r?\n)?\*\/`)
	match := re.FindAllStringSubmatch(string(data), -1)
	var title string
	var desc string
	var partners []string
	var depends []string
	if match == nil {
		return r, nil
	}
	title = strings.Trim(match[0][1], " \r\n")
	desc = strings.Trim(match[0][2], " \r\n")
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
		Title:    title,
		Desc:     desc,
		Partners: partners,
		Depends:  depends,
	}
	return r, nil
}

// trimAll takes the elements of a slice of strings and trims all the whitespace off the strings in the slice
func trimAll(input []string) []string {
	output := make([]string, len(input))
	for i, s := range input {
		output[i] = strings.Trim(s, " \r\n")
	}
	return output
}
