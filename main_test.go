package main

import (
	"io/fs"
	"io/ioutil"
	"testing"
)

func Test_compareFiles(t *testing.T) {
	cases := []struct {
		before     string
		after      string
		path       string
		breakpoint int
		diff       bool
	}{
		{
			before:     "./testdata/sample.png",
			after:      "./testdata/sample_copy.png",
			path:       "sample",
			breakpoint: 768,
			diff:       false,
		},
		{
			before:     "./testdata/sample.png",
			after:      "./testdata/sample_diff.png",
			path:       "sample",
			breakpoint: 768,
			diff:       true,
		},
	}

	for _, tt := range cases {
		beforelen := countFiles("./results")
		compareFiles(tt.before, tt.after, tt.path, tt.breakpoint)
		afterlen := countFiles(("./results"))

		if tt.diff {
			if beforelen == afterlen {
				t.Error(beforelen, afterlen)
			}
		} else {
			if beforelen != afterlen {
				t.Error(beforelen, afterlen)
			}
		}
	}
}

func countFiles(dirpath string) int {
	files, _ := ioutil.ReadDir(dirpath)
	var filesInfos []fs.FileInfo
	for _, fi := range files {
		filesInfos = append(filesInfos, fi)
	}
	return len(filesInfos)
}
