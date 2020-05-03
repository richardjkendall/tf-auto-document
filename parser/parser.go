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
	"github.com/zclconf/go-cty/cty/gocty"
)

// Parser is the object which holds the methods needed to scan hcl files looking for the details we need
type Parser struct {
	hclParser *hclparse.Parser
}

// ModuleDetails contains the details of the module being scanned
type ModuleDetails struct {
	title     string
	desc      string
	partners  []string
	depends   []string
	variables []VariableDetails
	outputs   []OutputDetails
}

// VariableDetails contains the details of the variables defined by the module
type VariableDetails struct {
	name     string
	desc     string
	def      string
	dataType string
}

// OutputDetails contains the details of the outputs defined by the module
type OutputDetails struct {
	name string
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
		varDetails.name = variableName
		fmt.Printf("\t\tvariable name: %s\n", variableName)
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
				varDetails.dataType = typeexpr.TypeString(valType)
			}
			// get description
			if attribute.Name == "description" && val.Type() == cty.String {
				varDetails.desc = val.AsString()
			}
			// get default
			if attribute.Name == "default" {
				// deal with string version
				if val.Type() == cty.String {
					varDetails.def = val.AsString()
				}
				// deal with numeric version
				if val.Type() == cty.Number {
					varDetails.def = val.AsBigFloat().String()
				}
				// deal with boolean version
				if val.Type() == cty.Bool {
					var ret bool
					gocty.FromCtyValue(val, &ret)
					if ret {
						varDetails.def = "true"
					} else {
						varDetails.def = "false"
					}
				}
				// deal with tuples and lists
				//...
			}
			// convert tuples
			if cty.Type.IsTupleType(val.Type()) {
				l := parser.convertTupleOrList(val)
				fmt.Printf("\t\t\t %s is a tuple: %v\n", attribute.Name, l)
			}
		}
		v = append(v, varDetails)
	}
	r.variables = v

	fmt.Printf("%+v\n", r)
	return r, nil
}

func (parser *Parser) convertTupleOrList(val cty.Value) []string {
	var ret []string

	for it := val.ElementIterator(); it.Next(); {
		_, v := it.Element()
		ret = append(ret, v.AsString())
	}
	return ret
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
	var re = regexp.MustCompile(`(?m)^\/\*\ntitle:\s+([\w\-]+)\ndesc:\s+([\w\-\s\.,/]+)\n(partners:\s+[\w\-,\s]+\n)?(depends:\s+[\w\-,\s]+\n)?\*\/`)
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
		output[i] = strings.Trim(s, " \r\n")
	}
	return output
}
