package validators

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

type InValidator struct {
	inMap      map[string]bool
	constraint string
}

func (vld *InValidator) ValidateValue(v interface{}) error {
	valid := true
	switch vt := v.(type) {
	case int:
		key := strconv.Itoa(vt)
		if _, found := vld.inMap[key]; !found {
			valid = false
		}
	case string:
		if _, found := vld.inMap[vt]; !found {
			valid = false
		}
	default:
		return ErrInvalidType

	}
	if !valid {
		return errors.Errorf(
			"input value '%v' is not in the '%s' set",
			v, vld.constraint,
		)
	}

	return nil
}

func (vld *InValidator) GetValidKinds() []reflect.Kind {
	return []reflect.Kind{reflect.Int, reflect.String}
}

func (vld *InValidator) SetConstraint(c string) error {
	if len(c) == 0 {
		return ErrInvalidConstraintValue
	}

	parts := strings.Split(c, ",")
	vld.inMap = make(map[string]bool, len(parts))
	for _, str := range parts {
		vld.inMap[str] = true
	}

	vld.constraint = c

	return nil
}
