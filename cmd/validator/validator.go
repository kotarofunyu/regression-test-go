package validator

import (
	"fmt"
	"strings"
)

func ValidateUrl(url, name string) error {
	if !strings.HasSuffix(url, "/") {
		return fmt.Errorf("%s must end with '/'", name)
	}
	return nil
}
