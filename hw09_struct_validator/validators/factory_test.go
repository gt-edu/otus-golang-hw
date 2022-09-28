package validators

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestFactory(t *testing.T) {
	type args struct {
		name        string
		fieldName   string
		structField reflect.StructField
		constraint  string
	}

	var tests = []struct {
		name          string
		args          args
		wantValidator Validator
		wantErr       error
	}{
		{
			name:          "fake",
			args:          args{name: "fake", fieldName: "fakeField", structField: reflect.StructField{}, constraint: "12"},
			wantValidator: nil,
			wantErr:       ErrUnavailableValidator,
		},
		{
			name: "max not int",
			args: args{name: "max", fieldName: "age", structField: reflect.ValueOf(struct {
				Age int `max:"2"`
			}{0}).Type().Field(0), constraint: "1{2"},
			wantValidator: nil,
			wantErr:       ErrInvalidConstraintValue,
		},
		{
			name: "max have to be int",
			args: args{name: "max", fieldName: "age", structField: reflect.ValueOf(struct {
				Age string `max:"2"`
			}{""}).Type().Field(0), constraint: "12"},
			wantValidator: nil,
			wantErr:       ErrInvalidType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValidator, err := Factory(tt.args.name, tt.args.structField, tt.args.constraint)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.wantValidator, gotValidator)
		})
	}
}
