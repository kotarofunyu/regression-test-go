package main

import (
	"io/ioutil"
	"testing"
)

func Test_compareFiles(t *testing.T) {
	before := "./testdata/sample.png"
	after := "./testdata/sample_copy.png"
	path := "./testdata"
	breakpoint := 768
	compareFiles(before, after, path, breakpoint)
	all, _ := ioutil.ReadDir("./")
	var exits bool
	for _, e := range all {
		if e.Name() == "results" {
			exits = true
			return
		}
	}

	if exits != true {
		t.Errorf("error: %#v", exits)
	}
}
