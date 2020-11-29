package manta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/influxdata/cron"
)

var (
	ErrFieldMustBeSet = errors.New("failed must be set")
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

func (m *Check) Validate() error {
	if m.Name == "" {
		return invalidField("name", ErrFieldMustBeSet)
	}

	if m.Desc == "" {
		return invalidField("desc", ErrFieldMustBeSet)
	}

	if m.Expr == "" {
		return invalidField("expr", ErrFieldMustBeSet)
	}

	if m.Status == "" {
		return invalidField("status", ErrFieldMustBeSet)
	}

	if err := validateStatus(m.Status); err != nil {
		return invalidField("status", err)
	}

	if _, err := cron.ParseUTC(m.Cron); err != nil {
		return invalidField("cron", err)
	}

	if m.Datasource == 0 {
		return invalidField("datasource", ErrFieldMustBeSet)
	}

	// todo: validate conditions

	return nil
}

func (m *Matcher) Validate() error {
	if m.Type < 0 && m.Type > 3 {
		return invalidField("matcher.type")
	}

	if m.Name == "" {
		return invalidField("matcher.Name", ErrFieldMustBeSet)
	}

	if m.Value == "" {
		return invalidField("matchers.value", ErrFieldMustBeSet)
	}

	return nil
}

func (m *Template) Validate() error {
	if m.Name != "" {
		return invalidField("name", ErrFieldMustBeSet)
	}

	if m.Desc != "" {
		return invalidField("desc", ErrFieldMustBeSet)
	}

	if len(m.Matchers) == 0 {
		return invalidField("matchers", ErrFieldMustBeSet)
	}

	for _, matcher := range m.Matchers {
		if err := matcher.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func inStrings(s string, slices []string) bool {
	for i := 0; i < len(slices); i++ {
		if slices[i] == s {
			return true
		}
	}

	return false
}

func (m *Threshold) Validate() error {
	if !inStrings(m.Type, thresholdTypes) {
		return invalidField("status", fmt.Errorf("invalid severity %q", m.Type))
	}

	if m.Type == Inside || m.Type == Outside {
		if m.Min > m.Max {
			return fmt.Errorf("threshold's min can't be larger than max")
		}
	}

	return nil
}
