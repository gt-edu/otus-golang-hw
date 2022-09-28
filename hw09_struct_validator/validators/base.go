package validators

import (
	"github.com/pkg/errors"
	"reflect"
)

var (
	ErrValueIsNotStruct       = errors.New("value is not a struct")
	ErrInvalidValidateTag     = errors.New("validate tag contains invalid value")
	ErrUnavailableValidator   = errors.New("validate tag contains unavailable validator")
	ErrInvalidType            = errors.New("field is of invalid type")
	ErrInvalidConstraintValue = errors.New("constraint value is invalid")
)

type Validator interface {
	HasValidType(rfType reflect.Type) bool
	ValidateValue(v interface{}) error
	SetConstraint(c string)
	ParseConstraint() error
}

type BaseValidator struct {
	constraint string
	kind       reflect.Kind
}

func (vld *BaseValidator) SetConstraint(c string) {
	vld.constraint = c
}
