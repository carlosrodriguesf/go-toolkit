package pipeline

import "fmt"

type RunError struct {
	StepName string
	Err      error
}

func (r RunError) Error() string {
	return fmt.Sprintf("error running step '%s': %s", r.StepName, r.Err)
}

func (r RunError) Unwrap() error {
	return r.Err
}
