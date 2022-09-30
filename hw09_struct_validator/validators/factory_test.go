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
		{
			name: "get max",
			args: args{name: "max", fieldName: "age", structField: reflect.ValueOf(struct {
				Age int `validate:"max:2"`
			}{0}).Type().Field(0), constraint: "2"},
			wantValidator: &MaxValidator{},
			wantErr:       nil,
		},
		{
			name: "get min",
			args: args{name: "min", fieldName: "age", structField: reflect.ValueOf(struct {
				Age int `validate:"min:2"`
			}{0}).Type().Field(0), constraint: "2"},
			wantValidator: &MinValidator{},
			wantErr:       nil,
		},
		{
			name: "get len",
			args: args{name: "len", fieldName: "name", structField: reflect.ValueOf(struct {
				Name string `validate:"len:4"`
			}{"John"}).Type().Field(0), constraint: "4"},
			wantValidator: &LenValidator{},
			wantErr:       nil,
		},
		{
			name: "get in",
			args: args{name: "in", fieldName: "age", structField: reflect.ValueOf(struct {
				Age int `validate:"in:2,3"`
			}{2}).Type().Field(0), constraint: "2,3"},
			wantValidator: &InValidator{},
			wantErr:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValidator, err := Factory(tt.args.name, tt.args.structField, tt.args.constraint)
			require.ErrorIs(t, err, tt.wantErr)
			require.IsType(t, tt.wantValidator, gotValidator)
		})
	}
}
