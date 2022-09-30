package validators

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInValidator_ValidateValue(t *testing.T) {
	validator := InValidator{}
	err := validator.SetConstraint("-1,0,1,2,3")
	require.NoError(t, err)

	tests := []struct {
		v      interface{}
		errMsg string
	}{
		{v: -1, errMsg: ""},
		{v: 0, errMsg: ""},
		{v: 1, errMsg: ""},
		{v: 3, errMsg: ""},
		{v: -4, errMsg: "input value '-4' is not in the '-1,0,1,2,3' set"},
		{v: "-1", errMsg: ""},
		{v: "0", errMsg: ""},
		{v: "1", errMsg: ""},
		{v: "3", errMsg: ""},
		{v: "-4.1", errMsg: "input value '-4.1' is not in the '-1,0,1,2,3' set"},
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

func TestInValidator_SetConstraint(t *testing.T) {
	tests := []struct {
		name    string
		c       string
		wantErr error
	}{
		{"valid case", "1,2", nil},
		{"invalid case - empty", "", ErrInvalidConstraintValue},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vld := &InValidator{}
			err := vld.SetConstraint(tt.c)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
