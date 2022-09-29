package validators

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func testIntSetConstraint(t *testing.T, vld Validator, negativeOk bool) {
	type testf struct {
		name    string
		c       string
		wantErr error
	}
	tests := []testf{
		{"simple case", "2", nil},
		{"zero", "0", nil},
		{"incorrect int", "1{2", ErrInvalidConstraintValue},
		{"whitespaces", " ", ErrInvalidConstraintValue},
	}

	if negativeOk {
		tests = append(tests, testf{"negative", "-1", nil})
	} else {
		tests = append(tests, testf{"negative", "-1", ErrInvalidConstraintValue})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vld.SetConstraint(tt.c)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
