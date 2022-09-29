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
	GetValidKinds() []reflect.Kind
	ValidateValue(v interface{}) error
	SetConstraint(c string) error
}
