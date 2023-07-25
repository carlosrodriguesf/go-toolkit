package groupederror_test

import (
	"encoding/json"
	"errors"
	"github.com/carlosrodriguesf/go-toolkit/groupederror"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestErrorf(t *testing.T) {
	var (
		err1           = errors.New("error 1")
		err2           = errors.New("error 2")
		errInvalidJSON = json.Unmarshal([]byte("invalid json"), &map[string]any{})
		errOut         = errors.New("not included")
	)

	err := groupederror.Errorf("%s: %s: %s", err1, err2, errInvalidJSON)

	require.EqualError(t, err, "error 1: error 2: invalid character 'i' looking for beginning of value")
	require.ErrorIs(t, err, err1)
	require.ErrorIs(t, err, err2)
	require.NotErrorIs(t, err, errOut)

	var syntaxError *json.SyntaxError
	require.ErrorAs(t, err, &syntaxError)
}
