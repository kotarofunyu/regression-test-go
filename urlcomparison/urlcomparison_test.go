package urlcomparison

import "testing"

func Test_newFileName(t *testing.T) {
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
	n := newFileName("before", "hogehoge", 1024)
	if n != "./captures/before-hogehoge-1024.png" {
		t.Error("error")
	}

	for _, test := range tests {
		got := newFileName(test.timing, test.path, test.breakpoint)
		if got != test.want {
			t.Errorf("expected %s, but got %s", test.want, got)
		}
	}
}
