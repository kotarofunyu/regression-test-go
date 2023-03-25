package urlcomparison

import (
	"os"
	"testing"
)

func Test_createOutputDir(t *testing.T) {
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
		createOutputDir(test.resultDir, test.captureDir)
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
