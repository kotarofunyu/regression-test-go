package comparison

import (
	"reflect"
	"testing"
)

func TestSetupBrowser(t *testing.T) {
	p, d := SetupBrowser()
	p_want := "*agouti.Page"
	p_got := reflect.TypeOf(p).String()
	d_want := "*agouti.WebDriver"
	d_got := reflect.TypeOf(d).String()

	if p_want != p_got {
		t.Error("error", p_got)
	}
	if d_want != d_got {
		t.Error("error", d_got)
	}
}
