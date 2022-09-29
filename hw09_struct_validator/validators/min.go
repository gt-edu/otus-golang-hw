package validators

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

type MinValidator struct {
	min int
}

func (vld *MinValidator) ValidateValue(v interface{}) error {
	valid := true
	switch vt := v.(type) {
	case int:
		if vt < vld.min {
			valid = false
		}
	case float32:
		if vt < float32(vld.min) {
			valid = false
		}
	case float64:
		if vt < float64(vld.min) {
			valid = false
		}
	default:
		return ErrInvalidType

	}
	if !valid {
		return errors.Errorf(
			"input value '%v' less then minimum '%d'",
			v, vld.min,
		)
	}

	return nil
}

func (vld *MinValidator) GetValidKinds() []reflect.Kind {
	return []reflect.Kind{reflect.Int, reflect.Float64, reflect.Float32}
}

func (vld *MinValidator) SetConstraint(c string) error {
	min, err := strconv.Atoi(c)
	if err != nil {
		return errors.Wrap(ErrInvalidConstraintValue, err.Error())
	}

	vld.min = min

	return nil
}
