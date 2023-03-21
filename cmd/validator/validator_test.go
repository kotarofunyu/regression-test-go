package validator

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestValidateUrl(t *testing.T) {
	got := ValidateUrl("http://localhost/", "url")
	if got != nil {
		t.Errorf("expected %#v, but got %s", nil, got)
	}

	got = ValidateUrl("http://localhost", "url")
	want := reflect.TypeOf(errors.New("hoge")).String()
	fmt.Println(want)
	if reflect.TypeOf(got).String() != want {
		t.Errorf("expected %s, but got %s", want, got)
	}
}
