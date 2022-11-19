package multierr

import (
	"fmt"
	"strings"
)

type Errors []error

func (es *Errors) Add(err error) {
	if err == nil {
		return
	}

	*es = append(*es, err)
}

func (es *Errors) Unwrap() error {
	if es == nil {
		return nil
	}

	if len(*es) == 0 {
		return nil
	}

	var c chain
	for _, e := range *es {
		c = append(c, e)
	}

	return c
}

// Error implements the error interface
func (es *Errors) Error() string {
	return es.Unwrap().Error()
}

// chain implements the interfaces necessary for errors.Is/As/Unwrap to
// work in a deterministic way with multierror. A chain tracks a list of
// errors while accounting for the current represented error. This lets
// Is/As be meaningful.
//
// Unwrap returns the next error. In the cleanest form, Unwrap would return
// the wrapped error here but we can't do that if we want to properly
// get access to all the errors. Instead, users are recommended to use
// Is/As to get the correct error type out.
//
// Precondition: []error is non-empty (len > 0)
type chain []error

// Error implements the error interface
func (e chain) Error() string {
	return listFormatFunc(e)
}

// listFormatFunc is a basic formatter that outputs the number of errors
// that occurred along with a bullet point list of the errors.
func listFormatFunc(es []error) string {
	if len(es) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %s\n\n", es[0])
	}

	points := make([]string, len(es))
	for i, err := range es {
		points[i] = fmt.Sprintf("* %s", err)
	}

	return fmt.Sprintf(
		"%d errors occurred:\n\t%s\n\n",
		len(es), strings.Join(points, "\n\t"))
}
