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
	valid := true
	switch vt := v.(type) {
	case int:
		if vt > vld.max {
			valid = false
		}
	case float32:
		if vt > float32(vld.max) {
			valid = false
		}
	case float64:
		if vt > float64(vld.max) {
			valid = false
		}
	default:
		return ErrInvalidType

	}
	if !valid {
		return errors.Errorf(
			"input value '%v' is greater then maximum '%d'",
			v, vld.max,
		)
	}

	return nil
}

func (vld *MaxValidator) GetValidKinds() []reflect.Kind {
	return []reflect.Kind{reflect.Int, reflect.Float64, reflect.Float32}
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
