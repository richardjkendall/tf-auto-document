package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

func convertValueToString(val cty.Value) string {
	// for basic types we can return the string representation right away
	if val.Type() == cty.String {
		return val.AsString()
	}
	if val.Type() == cty.Number {
		return numberToString(val)
	}
	if val.Type() == cty.Bool {
		return boolToString(val)
	}
	// for tuples (which seems to include lists) we need to iterate
	if val.Type().IsTupleType() {
		var ret []string
		for it := val.ElementIterator(); it.Next(); {
			_, v := it.Element()
			ret = append(ret, convertValueToString(v))
		}
		return "[" + strings.Join(ret, ", ") + "]"
	}
	// for objects (which seems to include maps) we need to iterate though the attributes
	// need to sort attributes first so that comparisons will work
	if val.Type().IsObjectType() {
		var ret []string
		atys := val.Type().AttributeTypes()
		attributeNames := make([]string, 0, len(atys))
		for name := range atys {
			attributeNames = append(attributeNames, name)
		}
		sort.Strings(attributeNames)
		for _, name := range attributeNames {
			ret = append(ret, name+"="+convertValueToString(val.GetAttr(name)))
		}
		return "{" + strings.Join(ret, ", ") + "}"
	}
	// if we get here we have an issue
	return "ERROR: cannot convert!"
}

func convertMap(val cty.Value) []string {
	var ret []string
	fmt.Printf("in map convert\n")

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
