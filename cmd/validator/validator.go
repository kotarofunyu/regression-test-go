package validator

import (
	"errors"
	"strings"
)

func ValidateUrl(url, name string) error {
	if !strings.HasSuffix(url, "/") {
		return errors.New(name + " must end with '/'")
	}
	return nil
}
