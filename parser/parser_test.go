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
