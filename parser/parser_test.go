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
