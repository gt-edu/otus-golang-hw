package validators

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLenValidator_ValidateValue(t *testing.T) {
	validator := LenValidator{}
	err := validator.SetConstraint("10")
	require.NoError(t, err)

	tests := []struct {
		v      interface{}
		errMsg string
	}{
		{v: "1234567890", errMsg: ""},
		{v: "123456789界", errMsg: ""},
		{v: "12345678界", errMsg: "input value '12345678界' has a length not equal '10'"},
		{v: "", errMsg: "input value '' has a length not equal '10'"},
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

func TestLenValidator_SetConstraint(t *testing.T) {
	vld := &LenValidator{}

	testIntSetConstraint(t, vld, false)
}
