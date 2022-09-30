package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gt-edu/otus-golang-hw/hw09_struct_validator/validators"
	"github.com/stretchr/testify/require"
)

type UserRole string

type CustomUser User

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		//{
		//	in: User{
		//		ID:     "1",
		//		Name:   "Test",
		//		Age:    22,
		//		Email:  "test@example.org",
		//		Role:   "admin",
		//		Phones: []string{"12345678901"},
		//	},
		//	expectedErr: nil,
		//},
		//{
		//	in: CustomUser{
		//		ID:     "1",
		//		Name:   "Test",
		//		Age:    22,
		//		Email:  "test@example.org",
		//		Role:   "admin",
		//		Phones: []string{"12345678901"},
		//	},
		//	expectedErr: nil,
		//},
		{
			name:        "value is not struct",
			in:          UserRole("test"),
			expectedErr: validators.ErrValueIsNotStruct,
		},
		{
			name: "max-min simple case",
			in: struct {
				Age  int     `validate:"max:2|min:0"`
				AgeF float32 `validate:"max:2|min:0"`
			}{1, float32(1.0)},
			expectedErr: nil,
		},
		{
			name: "max-min simple case slices",
			in: struct {
				Ages  []int     `validate:"max:2|min:0"`
				AgesF []float64 `validate:"max:2|min:0"`
			}{[]int{1, 2}, []float64{1.0, 2.0}},
			expectedErr: nil,
		},
		{
			name: "max-min value bigger then",
			in: struct {
				Age  int     `validate:"max:2|min:0"`
				AgeF float32 `validate:"max:2|min:0"`
			}{3, float32(3.1)},
			expectedErr: ValidationErrors{
				NewValidationError("Age", "input value '3' is greater then maximum '2'"),
				NewValidationError("AgeF", "input value '3.1' is greater then maximum '2'"),
			},
		},
		{
			name: "max-min slice value bigger then",
			in: struct {
				Ages  []int     `validate:"max:2|min:0"`
				AgesF []float64 `validate:"max:2|min:0"`
			}{[]int{3, 2}, []float64{3.1, 2.0}},
			expectedErr: ValidationErrors{
				NewValidationError("Ages", "input value '3' is greater then maximum '2'"),
				NewValidationError("AgesF", "input value '3.1' is greater then maximum '2'"),
			},
		},
		{
			name: "max-min value less then",
			in: struct {
				Age  int     `validate:"max:2|min:0"`
				AgeF float32 `validate:"max:2|min:0"`
			}{-1, float32(-1.1)},
			expectedErr: ValidationErrors{
				NewValidationError("Age", "input value '-1' less then minimum '0'"),
				NewValidationError("AgeF", "input value '-1.1' less then minimum '0'"),
			},
		},
		{
			name: "max-min slice value less then",
			in: struct {
				Ages  []int     `validate:"max:2|min:0"`
				AgesF []float64 `validate:"max:2|min:0"`
			}{[]int{-1, 2}, []float64{-1.1, 2.0}},
			expectedErr: ValidationErrors{
				NewValidationError("Ages", "input value '-1' less then minimum '0'"),
				NewValidationError("AgesF", "input value '-1.1' less then minimum '0'"),
			},
		},
		{
			name: "len simple case",
			in: struct {
				Phone string `validate:"len:10"`
			}{"123456789界"},
			expectedErr: nil,
		},
		{
			name: "len simple case slices",
			in: struct {
				Phones []string `validate:"len:10"`
			}{[]string{"123456789界", "000456789界"}},
			expectedErr: nil,
		},
		{
			name: "len note equal",
			in: struct {
				Phone string `validate:"len:10"`
			}{"12345678界"},
			expectedErr: ValidationErrors{
				NewValidationError("Phone", "input value '12345678界' has a length not equal '10'"),
			},
		},
		{
			name: "len note equal slices",
			in: struct {
				Phones []string `validate:"len:10"`
			}{[]string{"12345678界", "00045678界"}},
			expectedErr: ValidationErrors{
				NewValidationError("Phones", "input value '12345678界' has a length not equal '10'"),
				NewValidationError("Phones", "input value '00045678界' has a length not equal '10'"),
			},
		},
		{
			name: "in simple case",
			in: struct {
				Phone  string `validate:"in:123,456"`
				PhoneI int    `validate:"in:123,456"`
			}{"456", 123},
			expectedErr: nil,
		},
		{
			name: "in simple case - slices",
			in: struct {
				Phone  []string `validate:"in:123,456"`
				PhoneI []int    `validate:"in:123,456"`
			}{[]string{"456", "123"}, []int{123, 456}},
			expectedErr: nil,
		},
		{
			name: "in simple case invalid",
			in: struct {
				Phone  string `validate:"in:123,456"`
				PhoneI int    `validate:"in:123,456"`
			}{"4567", 1234},
			expectedErr: ValidationErrors{
				NewValidationError("Phone", "input value '4567' is not in the '123,456' set"),
				NewValidationError("PhoneI", "input value '1234' is not in the '123,456' set"),
			},
		},
		{
			name: "in simple case invalid - slices",
			in: struct {
				Phone  []string `validate:"in:123,456"`
				PhoneI []int    `validate:"in:123,456"`
			}{[]string{"4567", "123"}, []int{1234, 456}},
			expectedErr: ValidationErrors{
				NewValidationError("Phone", "input value '4567' is not in the '123,456' set"),
				NewValidationError("PhoneI", "input value '1234' is not in the '123,456' set"),
			},
		},
		{
			name: "regexp simple case",
			in: struct {
				Phone string `validate:"regexp:^\\d+$"`
			}{"456"},
			expectedErr: nil,
		},
		{
			name: "regexp simple case - slices",
			in: struct {
				Phone []string `validate:"regexp:^\\d+$"`
			}{[]string{"456", "123"}},
			expectedErr: nil,
		},
		{
			name: "regexp simple case invalid",
			in: struct {
				Phone string `validate:"regexp:^\\d+$"`
			}{"456z"},
			expectedErr: ValidationErrors{
				NewValidationError("Phone", "input value '456z' does not match the '^\\d+$' regex"),
			},
		},
		{
			name: "in simple case invalid - slices",
			in: struct {
				Phone []string `validate:"regexp:^\\d+$"`
			}{[]string{"456z", "123z", "987"}},
			expectedErr: ValidationErrors{
				NewValidationError("Phone", "input value '456z' does not match the '^\\d+$' regex"),
				NewValidationError("Phone", "input value '123z' does not match the '^\\d+$' regex"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.NotNil(t, err)
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			}

			// Place your code here.
			_ = tt
		})
	}
}

func NewValidationError(fieldName, errMsg string) ValidationError {
	return ValidationError{
		Field: fieldName,
		Err:   errors.New(errMsg),
	}
}

func TestValidationErrors_Error(t *testing.T) {
	tests := []struct {
		name string
		v    ValidationErrors
		want string
	}{
		{
			name: "simple case",
			v: ValidationErrors{
				ValidationError{Field: "f1", Err: errors.New("err1")},
				ValidationError{Field: "f2", Err: errors.New("err2")},
				ValidationError{Field: "f3", Err: errors.New("err3")},
			},
			want: "Validation errors: f1: err1; f2: err2; f3: err3",
		},
		{
			name: "empty",
			v:    ValidationErrors{},
			want: "Validation errors: ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.Error()
			require.Equal(t, tt.want, got)
		})
	}
}
