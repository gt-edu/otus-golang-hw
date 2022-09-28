package validators

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateValue(t *testing.T) {
	t.Run("test max validate value", func(t *testing.T) {
		validator := MaxValidator{max: 3}

		tests := []struct {
			v      int
			errMsg string
		}{
			{v: -1, errMsg: ""},
			{v: 0, errMsg: ""},
			{v: 1, errMsg: ""},
			{v: 3, errMsg: ""},
			{v: 4, errMsg: "input value '4' is greater then maximum '3'"},
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
