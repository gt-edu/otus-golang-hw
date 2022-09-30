package validators

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"unicode/utf8"
)

type LenValidator struct {
	len int
}

func (vld *LenValidator) ValidateValue(v interface{}) error {
	valid := utf8.RuneCountInString(fmt.Sprintf("%v", v)) == vld.len

	if !valid {
		return errors.Errorf(
			"input value '%v' has a length not equal '%d'",
			v, vld.len,
		)
	}

	return nil
}

func (vld *LenValidator) GetValidKinds() []reflect.Kind {
	return []reflect.Kind{reflect.String}
}

func (vld *LenValidator) SetConstraint(c string) error {
	length, err := strconv.Atoi(c)
	if err != nil {
		return errors.Wrap(ErrInvalidConstraintValue, err.Error())
	}

	if length < 0 {
		return ErrInvalidConstraintValue
	}

	vld.len = length

	return nil
}
