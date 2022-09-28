package validators

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

type MaxValidator struct {
	BaseValidator
	max int
}

func (vld *MaxValidator) ValidateValue(v interface{}) error {
	intValue := v.(int)

	if intValue > vld.max {
		return errors.Errorf(
			"input value '%d' is greater then maximum '%d'",
			intValue, vld.max,
		)
	}

	return nil
}

func (vld *MaxValidator) HasValidType(rfType reflect.Type) bool {
	if rfType.Kind() == reflect.Int {
		return true
	}

	return rfType.Kind() == reflect.Slice && rfType.Elem().Kind() == reflect.Int
}

func (vld *MaxValidator) ParseConstraint() error {
	max, err := strconv.Atoi(vld.constraint)
	if err != nil {
		return errors.Wrap(ErrInvalidConstraintValue, err.Error())
	}

	vld.max = max

	return nil
}
