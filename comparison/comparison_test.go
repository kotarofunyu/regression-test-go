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

func TestNewFileName(t *testing.T) {
	tests := []struct {
		timing     string
		path       string
		breakpoint int
		want       string
	}{
		{
			timing:     "before",
			path:       "hogehoge",
			breakpoint: 1024,
			want:       "./captures/before-hogehoge-1024.png",
		},
		{
			timing:     "after",
			path:       "fugafuga",
			breakpoint: 532,
			want:       "./captures/after-fugafuga-532.png",
		},
	}
	n := NewFileName("before", "hogehoge", 1024)
	if n != "./captures/before-hogehoge-1024.png" {
		t.Error("error")
	}

	for _, test := range tests {
		got := NewFileName(test.timing, test.path, test.breakpoint)
		if got != test.want {
			t.Errorf("expected %s, but got %s", test.want, got)
		}
	}
}
