package parser

import (
	"testing"

	"github.com/go-test/deep"
)

func TestMainDetails(t *testing.T) {
	want := ModuleDetails{
		Title:    "testing",
		Desc:     "test, test, test",
		Partners: []string{"partner1", "partner2"},
		Depends:  []string{"depend1", "depend2"},
	}
	got, err := New().getMainDetails("tests/main_test.tf")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}

func TestMainWithPartnersOnly(t *testing.T) {
	want := ModuleDetails{
		Title:    "testing",
		Desc:     "test, test, test",
		Partners: []string{"partner1", "partner2"},
	}
	got, err := New().getMainDetails("tests/main_test_p_only.tf")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}

func TestMainWithDependsOnly(t *testing.T) {
	want := ModuleDetails{
		Title:   "testing",
		Desc:    "test, test, test",
		Depends: []string{"depend1", "depend2"},
	}
	got, err := New().getMainDetails("tests/main_test_d_only.tf")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}

func TestMainWithPunc(t *testing.T) {
	want := ModuleDetails{
		Title: "testing",
		Desc:  `test with lots of punctuation <>=""':;!@#$%^&*()-_+*~.`,
	}
	got, err := New().getMainDetails("tests/main_test_punc.tf")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}

func TestSimpleVariable(t *testing.T) {
	want := ModuleDetails{
		Variables: []VariableDetails{
			VariableDetails{
				Name: "test",
				Desc: "testing variable with no type",
			},
		},
	}
	got, err := New().ParseFolder("tests/simple_variable/")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}

func TestTypedVariables(t *testing.T) {
	want := ModuleDetails{
		Variables: []VariableDetails{
			VariableDetails{
				Name:     "test_string",
				Desc:     "this is a string",
				DataType: "string",
				Def:      "string",
			},
			VariableDetails{
				Name:     "test_number",
				Desc:     "this is a number",
				DataType: "number",
				Def:      "10",
			},
			VariableDetails{
				Name:     "test_bool",
				Desc:     "this is a bool",
				DataType: "bool",
				Def:      "true",
			},
		},
	}
	got, err := New().ParseFolder("tests/variable_typed/")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}

func TestComplexVariables(t *testing.T) {
	want := ModuleDetails{
		Variables: []VariableDetails{
			VariableDetails{
				Name:     "test_string_list",
				Desc:     "list of strings",
				DataType: "list(string)",
				Def:      "[one, two, three]",
			},
			VariableDetails{
				Name:     "test_number_list",
				Desc:     "list of numbers",
				DataType: "list(number)",
				Def:      "[1, 2, 3]",
			},
			VariableDetails{
				Name:     "test_bool_list",
				Desc:     "list of bools",
				DataType: "list(bool)",
				Def:      "[true, false]",
			},
			VariableDetails{
				Name:     "test_tuple_mv",
				Desc:     "multi-value tuple",
				DataType: "tuple([string,number,bool])",
				Def:      "[test, 1, true]",
			},
			VariableDetails{
				Name:     "test_string_map",
				Desc:     "test map for strings",
				DataType: "map(string)",
				Def:      "{a=ay, b=bee, c=cee}",
			},
			VariableDetails{
				Name:     "test_object",
				Desc:     "test object",
				DataType: "object({a=string,b=number,c=bool})",
				Def:      "{a=ay, b=10, c=false}",
			},
			VariableDetails{
				Name:     "test_string_set",
				Desc:     "set of strings",
				DataType: "set(string)",
				Def:      "[one, two, three]",
			},
			VariableDetails{
				Name:     "test_list_of_objects",
				Desc:     "test list of objects",
				DataType: "list(object({a=string,b=number,c=bool}))",
				Def:      "[{a=ay, b=10, c=false}, {d=dee, e=20, f=true}]",
			},
			VariableDetails{
				Name:     "test_object_with_list",
				Desc:     "test object with a list",
				DataType: "object({a=list(string)})",
				Def:      "{a=[a, b, c]}",
			},
		},
	}
	got, err := New().ParseFolder("tests/variable_complex/")
	if err != nil {
		t.Errorf("Issue %q", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error(diff)
	}
}
