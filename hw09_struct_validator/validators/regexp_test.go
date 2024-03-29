package validators

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegexpValidator_ValidateValue(t *testing.T) {
	validator := RegexpValidator{}
	err := validator.SetConstraint("^\\d+$")
	require.NoError(t, err)

	tests := []struct {
		v      interface{}
		errMsg string
	}{
		{v: "1234567890", errMsg: ""},
		// {v: "123456789界", errMsg: ""},
		{v: "12345678界", errMsg: "input value '12345678界' does not match the '^\\d+$' regex"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			vldErr := validator.ValidateValue(tt.v)

			if len(tt.errMsg) > 0 {
				require.Equal(t, tt.errMsg, vldErr.Error())
			} else {
				require.Nil(t, vldErr)
			}
		})
	}
}

func TestRegexpValidator_SetConstraint(t *testing.T) {
	vld := &RegexpValidator{}

	tests := []struct {
		name    string
		c       string
		wantErr error
	}{
		{"valid case", "\\d+", nil},
		{"invalid case - empty", "", ErrInvalidConstraintValue},
		{"invalid case - invalid regexp", "?=()", ErrInvalidConstraintValue},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vld.SetConstraint(tt.c)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
