package main

import (
	"flag"
	"io/fs"
	"io/ioutil"
	"reflect"
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

func Test_setupArgs(t *testing.T) {
	flag.CommandLine.Set("base_url", "http://localhost")
	flag.CommandLine.Set("paths", "hoge,fuga")
	flag.CommandLine.Set("breakpoints", "100,200,300")
	flag.CommandLine.Set("gitpath", "./projecta")
	flag.CommandLine.Set("beforebranch", "develop")
	flag.CommandLine.Set("afterbranch", "feature_a")
	flag.CommandLine.Set("beforeurl", "http://example.com")
	flag.CommandLine.Set("afterurl", "http://localhost")
	baseurl, paths, breakpoints, gitpath, beforebranch, afterbranch, beforeurl, afterurl := setupArgs()
	if baseurl != "http://localhost" {
		t.Error("expected http://localhost, but got", baseurl)
	}
	if !reflect.DeepEqual([]string{"hoge", "fuga"}, paths) {
		t.Error("expexted, but got", paths)
	}
	if !reflect.DeepEqual([]int{100, 200, 300}, breakpoints) {
		t.Error("expexted, but got", paths)
	}
	if gitpath != "./projecta" {
		t.Error("error")
	}
	if beforebranch != "develop" {
		t.Error("error")
	}
	if afterbranch != "feature_a" {
		t.Error("error")
	}
	if beforeurl != "http://example.com" {
		t.Error("error")
	}
	if afterurl != "http://localhost" {
		t.Error("error")
	}
}
