package manta

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrFieldMustBeSet = errors.New("field must be set")
)

func invalidField(field string, errs ...error) error {
	if len(errs) == 0 {
		return fmt.Errorf("invalid field " + field)
	}

	es := make([]string, len(errs))
	for i := 0; i < len(errs); i++ {
		es[i] = errs[i].Error()
	}

	return fmt.Errorf("invalid field " + field + ", err: " + strings.Join(es, ";"))
}

func validateStatus(v string) error {
	if v == "active" || v == "inactive" {
		return nil
	}

	return fmt.Errorf("status must be active or inactive")
}
