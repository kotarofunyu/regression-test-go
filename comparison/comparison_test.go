package comparison

import (
	"os"
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

func TestCreateOutputDir(t *testing.T) {
	tests := []struct {
		description string
		resultDir   string
		captureDir  string
		before      func()
		after       func()
	}{
		{
			description: "It can create dirs",
			resultDir:   "hogetest",
			captureDir:  "fugatest",
			before:      func() {},
			after: func() {
				os.Remove("hogetest")
				os.Remove("fugatest")
			},
		},
		{
			description: "No Problem with trying to create existing dir",
			resultDir:   "hogetest",
			captureDir:  "fugatest",
			before: func() {
				os.Create("hogetest")
			},
			after: func() {
				os.Remove("hogetest")
				os.Remove("fugatest")
			},
		},
	}

	for _, test := range tests {
		test.before()
		CreateOutputDir(test.resultDir, test.captureDir)
		defer test.after()
		entries, err := os.ReadDir("./")
		var firstexists bool
		var secondexists bool
		if err != nil {
			t.Errorf("error running test %s", test.description)
		}
		for _, e := range entries {
			if e.Name() == test.resultDir {
				firstexists = true
			}
			if e.Name() == test.captureDir {
				secondexists = true
			}
		}

		if !(firstexists && secondexists) {
			t.Errorf("error running test %s", test.description)
		}
	}
}
