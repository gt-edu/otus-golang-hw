package validators

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateValue(t *testing.T) {
	t.Run("test max validate value", func(t *testing.T) {
		validator := MaxValidator{}
		err := validator.SetConstraint("3")
		require.NoError(t, err)

		tests := []struct {
			v      interface{}
			errMsg string
		}{
			{v: -1, errMsg: ""},
			{v: 0, errMsg: ""},
			{v: 1, errMsg: ""},
			{v: 3, errMsg: ""},
			{v: 4, errMsg: "input value '4' is greater then maximum '3'"},
			{v: -1.0, errMsg: ""},
			{v: 0.0, errMsg: ""},
			{v: 1.0, errMsg: ""},
			{v: 3.0, errMsg: ""},
			{v: 4.1, errMsg: "input value '4.1' is greater then maximum '3'"},
			{v: float32(-1.0), errMsg: ""},
			{v: float32(0.0), errMsg: ""},
			{v: float32(1.0), errMsg: ""},
			{v: float32(3.0), errMsg: ""},
			{v: float32(4.1), errMsg: "input value '4.1' is greater then maximum '3'"},
		}

		for i, tt := range tests {
			t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
				vldErr := validator.ValidateValue(tt.v)

				if len(tt.errMsg) > 0 {
					require.Equal(t, vldErr.Error(), tt.errMsg)
				} else {
					require.Nil(t, vldErr)
				}
			})
		}
	})
}

func TestMaxValidator_SetConstraint(t *testing.T) {
	tests := []struct {
		name    string
		c       string
		wantErr error
	}{
		{"simple case", "2", nil},
		{"negative", "-1", ErrInvalidConstraintValue},
		{"zero", "0", ErrInvalidConstraintValue},
		{"incorrect int", "1{2", ErrInvalidConstraintValue},
		{"whitespaces", " ", ErrInvalidConstraintValue},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vld := &MaxValidator{}
			err := vld.SetConstraint(tt.c)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
