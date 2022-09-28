package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
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
			in:          UserRole("test"),
			expectedErr: validators.ErrValueIsNotStruct,
		},
		{
			in: struct {
				Age int `validate:"max:2"`
			}{1},
			expectedErr: nil,
		},
		{
			in: struct {
				Age int `validate:"max:2"`
			}{3},
			expectedErr: ValidationErrors{
				NewValidationError("Age", "input value '3' is greater then maximum '2'"),
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, err.Error(), tt.expectedErr.Error())
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
