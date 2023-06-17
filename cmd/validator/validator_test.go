package validator

import (
	"testing"
)

func TestValidateUrl(t *testing.T) {
	cases := []struct {
		url  string
		name string
		want any
	}{
		{
			url:  "http://localhost",
			name: "url",
			want: "url must end with '/'",
		},
		{
			url:  "http://localhost/",
			name: "url",
			want: nil,
		},
	}
	for _, c := range cases {
		got := ValidateUrl(c.url, c.name)
		if c.want == nil {
			if got != nil {
				t.Errorf("expected %#v, but got %s", c.want, got)
			}
			continue
		}

		if got.Error() != c.want {
			t.Errorf("expected %#v, but got %s", c.want, got)
		}
	}
}
