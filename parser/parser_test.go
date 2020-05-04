package parser

import (
	"reflect"
	"testing"
)

func TestMainDetails(t *testing.T) {
	want := ModuleDetails{
		title:    "testing",
		desc:     "test, test, test",
		partners: []string{"partner1", "partner2"},
		depends:  []string{"depend1", "depend2"},
	}
	got, err := New().getMainDetails("tests/main_test.tf")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("getMainDetails() = %q, want %q", got, want)
	}
}

func TestMainWithPartnersOnly(t *testing.T) {
	want := ModuleDetails{
		title:    "testing",
		desc:     "test, test, test",
		partners: []string{"partner1", "partner2"},
	}
	got, err := New().getMainDetails("tests/main_test_p_only.tf")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("getMainDetails() = %q, want %q", got, want)
	}
}

func TestMainWithDependsOnly(t *testing.T) {
	want := ModuleDetails{
		title:   "testing",
		desc:    "test, test, test",
		depends: []string{"depend1", "depend2"},
	}
	got, err := New().getMainDetails("tests/main_test_d_only.tf")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("getMainDetails() = %q, want %q", got, want)
	}
}

func TestMainWithPunc(t *testing.T) {
	want := ModuleDetails{
		title: "testing",
		desc:  `test with lots of punctuation <>=""':;!@#$%^&*()-_+*~.`,
	}
	got, err := New().getMainDetails("tests/main_test_punc.tf")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("getMainDetails() = %q, want %q", got, want)
	}
}

func TestSimpleVariable(t *testing.T) {
	want := ModuleDetails{
		variables: []VariableDetails{
			VariableDetails{
				name: "test",
				desc: "testing variable with no type",
			},
		},
	}
	got, err := New().ParseFolder("tests/simple_variable/")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseFolder() = %q, want %q", got, want)
	}
}

func TestTypedVariables(t *testing.T) {
	want := ModuleDetails{
		variables: []VariableDetails{
			VariableDetails{
				name:     "test_string",
				desc:     "this is a string",
				dataType: "string",
				def:      "string",
			},
			VariableDetails{
				name:     "test_number",
				desc:     "this is a number",
				dataType: "number",
				def:      "10",
			},
			VariableDetails{
				name:     "test_bool",
				desc:     "this is a bool",
				dataType: "bool",
				def:      "true",
			},
		},
	}
	got, err := New().ParseFolder("tests/variable_typed/")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseFolder() = %q, want %q", got, want)
	}
}

func TestComplexVariables(t *testing.T) {
	want := ModuleDetails{
		variables: []VariableDetails{
			VariableDetails{
				name:     "test_string_list",
				desc:     "list of strings",
				dataType: "list(string)",
				def:      "[one, two, three]",
			},
			VariableDetails{
				name:     "test_number_list",
				desc:     "list of numbers",
				dataType: "list(number)",
				def:      "[1, 2, 3]",
			},
			VariableDetails{
				name:     "test_bool_list",
				desc:     "list of bools",
				dataType: "list(bool)",
				def:      "[true, false]",
			},
			VariableDetails{
				name:     "test_tuple_mv",
				desc:     "multi-value tuple",
				dataType: "tuple([string,number,bool])",
				def:      "[test, 1, true]",
			},
			VariableDetails{
				name:     "test_string_map",
				desc:     "test map for strings",
				dataType: "map(string)",
				def:      "b",
			},
		},
	}
	got, err := New().ParseFolder("tests/variable_complex/")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseFolder() = %q, want %q", got, want)
	}
}
