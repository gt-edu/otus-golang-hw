package validators

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

type MaxValidator struct {
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

	if rfType.Kind() == reflect.Slice && rfType.Elem().Kind() == reflect.Int {
		return true
	}

	return false
}

func (vld *MaxValidator) SetConstraint(c string) error {
	max, err := strconv.Atoi(c)
	if err != nil {
		return errors.Wrap(ErrInvalidConstraintValue, err.Error())
	}

	if max <= 0 {
		return errors.Wrap(ErrInvalidConstraintValue, fmt.Sprintf("'%d' less or equal zero", max))
	}

	vld.max = max

	return nil
}
