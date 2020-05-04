package parser

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func convertMap(val cty.Value) []string {
	var ret []string
	fmt.Printf("in map convert\n")

	return ret
}

func convertTupleOrList(val cty.Value) []string {
	var ret []string
	fmt.Printf("in convert tuple\n")
	// if list then all elements have the same type
	if val.Type().IsListType() {
		eleType := val.Type().ElementType()
		// handle strings
		if eleType == cty.String {
			for it := val.ElementIterator(); it.Next(); {
				_, v := it.Element()
				ret = append(ret, v.AsString())
			}
		}
	}
	// if tuple then elements can have different types
	if val.Type().IsTupleType() {
		eleTypes := val.Type().TupleElementTypes()
		var index int
		if !val.CanIterateElements() {
			fmt.Printf("can't iterate over elements\n")
		} else {
			fmt.Printf("can iterate over elements\n")
		}
		for it := val.ElementIterator(); it.Next(); {
			_, v := it.Element()
			fmt.Printf("type of element is %s\n", typeexpr.TypeString(v.Type()))
			// handle string
			if eleTypes[index] == cty.String {
				ret = append(ret, v.AsString())
			}
			// handle number
			if eleTypes[index] == cty.Number {
				ret = append(ret, numberToString(v))
			}
			// handle bool
			if eleTypes[index] == cty.Bool {
				ret = append(ret, boolToString(v))
			}

			index++
		}
	}
	return ret
}

func numberToString(val cty.Value) string {
	return val.AsBigFloat().String()
}

// boolToString converts cty.Value bool to a string representation
func boolToString(val cty.Value) string {
	var ret bool
	gocty.FromCtyValue(val, &ret)
	if ret {
		return "true"
	}
	return "false"
}
