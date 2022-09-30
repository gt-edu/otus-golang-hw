package validators

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type InValidator struct {
	inMap      map[string]bool
	constraint string
}

func (vld *InValidator) ValidateValue(v interface{}) error {
	valid := true
	switch reflect.ValueOf(v).Kind() {
	case reflect.Int:
		if _, found := vld.inMap[fmt.Sprintf("%d", v)]; !found {
			valid = false
		}
	case reflect.String:
		if _, found := vld.inMap[fmt.Sprintf("%v", v)]; !found {
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
