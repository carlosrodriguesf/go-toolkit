package groupederror

import (
	"errors"
	"fmt"
)

type GroupedError struct {
	template string
	errs     []error
}

func Errorf(template string, errs ...error) error {
	return GroupedError{
		template: template,
		errs:     errs,
	}
}

func (e GroupedError) Error() string {
	str := make([]any, len(e.errs))
	for i, err := range e.errs {
		str[i] = err.Error()
	}
	return fmt.Sprintf(e.template, str...)
}

func (e GroupedError) Is(err error) bool {
	for _, e := range e.errs {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}

func (e GroupedError) As(t any) bool {
	for _, e := range e.errs {
		if errors.As(e, t) {
			return true
		}
	}
	return false
}
