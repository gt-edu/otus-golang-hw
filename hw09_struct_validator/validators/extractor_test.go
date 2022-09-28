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

func TestExtractValidators(t *testing.T) {
	type args struct {
		rfVal reflect.Value
	}
	tests := []struct {
		name    string
		args    args
		want    [][]Validator
		wantErr error
	}{
		{
			name: "empty constraint",
			args: args{rfVal: reflect.ValueOf(struct {
				Age int `validate:"max:"`
			}{0})},
			want:    nil,
			wantErr: ErrInvalidValidateTag,
		},
		{
			name: "no colon",
			args: args{rfVal: reflect.ValueOf(struct {
				Age int `validate:"max"`
			}{0})},
			want:    nil,
			wantErr: ErrInvalidValidateTag,
		},
		{
			name: "only colon",
			args: args{rfVal: reflect.ValueOf(struct {
				Age int `validate:":"`
			}{0})},
			want:    nil,
			wantErr: ErrInvalidValidateTag,
		},
		{
			name: "invalid type",
			args: args{rfVal: reflect.ValueOf(struct {
				Age string `validate:"max:6"`
			}{""})},
			want:    nil,
			wantErr: ErrInvalidType,
		},
		{
			name: "unavailable validator",
			args: args{rfVal: reflect.ValueOf(struct {
				Age int `validate:"fake:z"`
			}{0})},
			want:    nil,
			wantErr: ErrUnavailableValidator,
		},
		{
			name:    "incorrect value type",
			args:    args{rfVal: reflect.ValueOf("")},
			want:    nil,
			wantErr: ErrValueIsNotStruct,
		},
		{
			name: "simple case",
			args: args{rfVal: reflect.ValueOf(struct {
				Age int `validate:"max:2"`
			}{0})},
			want:    [][]Validator{{&MaxValidator{max: 2}}},
			wantErr: nil,
		},
		{
			name: "many fields",
			args: args{rfVal: reflect.ValueOf(struct {
				Name  string
				Age   int `validate:"max:2"`
				Phone string
			}{"", 0, ""})},
			want:    [][]Validator{nil, {&MaxValidator{max: 2}}, nil},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractValidators(tt.args.rfVal)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.want == nil {
				require.Nil(t, tt.want)
			} else {
				require.Len(t, got, len(tt.want))
				for i, wantVldList := range tt.want {
					require.Len(t, got[i], len(tt.want[i]))
					for ii, wantVld := range wantVldList {
						require.IsType(t, reflect.TypeOf(wantVld), reflect.TypeOf(got[i][ii]))
					}
				}
			}
		})
	}
}
