package validators

import (
	"reflect"
	"regexp"

	"github.com/pkg/errors"
)

type RegexpValidator struct {
	rgxp       *regexp.Regexp
	constraint string
}

func (vld *RegexpValidator) ValidateValue(v interface{}) error {
	valid := true

	if !vld.rgxp.MatchString(v.(string)) {
		valid = false
	}

	if !valid {
		return errors.Errorf(
			"input value '%v' does not match the '%s' regex", v, vld.constraint,
		)
	}

	return nil
}

func (vld *RegexpValidator) GetValidKinds() []reflect.Kind {
	return []reflect.Kind{reflect.String}
}

func (vld *RegexpValidator) SetConstraint(c string) error {
	if len(c) == 0 {
		return ErrInvalidConstraintValue
	}

	rgxp, err := regexp.Compile(c)
	if err != nil {
		return errors.Wrap(ErrInvalidConstraintValue, err.Error())
	}

	vld.rgxp = rgxp
	vld.constraint = c

	return nil
}
