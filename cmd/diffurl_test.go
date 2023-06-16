package cmd

import (
	"bytes"
	"os"
	"strconv"
	"testing"
)

var originalArgs = os.Args

func Test_diffurlCmd(t *testing.T) {
	cases := []struct {
		beforeurl   string
		afterurl    string
		paths       []string
		breakpoints []int
		want        string
	}{
		{
			beforeurl:   "https://www.google.com/",
			afterurl:    "https://www.google.com/",
			paths:       []string{"sample"},
			breakpoints: []int{768},
			want:        "Image is same!",
		},
	}
	for _, tt := range cases {
		os.Args = append(os.Args, "diffurl",
			"--beforeurl", tt.beforeurl,
			"--afterurl", tt.afterurl,
			"--paths", tt.paths[0],
			"--breakpoints", strconv.Itoa(tt.breakpoints[0]),
		)
		defer func() { os.Args = originalArgs }()
		got := PickStdout(t, func() { diffurlCmd.Execute() })
		if got != tt.want {
			t.Errorf("subCmd.Execute() = %v, want = %v", got, tt.want)
		}
	}
}

func PickStdout(t *testing.T, fnc func()) string {
	t.Helper()
	backup := os.Stdout
	defer func() {
		os.Stdout = backup
	}()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("fail pipe: %v", err)
	}
	os.Stdout = w
	fnc()
	w.Close()
	var buffer bytes.Buffer
	if n, err := buffer.ReadFrom(r); err != nil {
		t.Fatalf("fail read buf: %v - number: %v", err, n)
	}
	s := buffer.String()
	return s[:len(s)-1]
}
