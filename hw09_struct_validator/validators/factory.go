package validators

import "reflect"

func Factory(name string, structField reflect.StructField, constraint string) (Validator, error) {
	var validator Validator = nil
	if name == "max" {
		validator = &MaxValidator{}
	}

	if validator == nil {
		return nil, ErrUnavailableValidator
	}

	if !validator.HasValidType(structField.Type) {
		return nil, ErrInvalidType
	}

	err := validator.SetConstraint(constraint)
	if err != nil {
		return nil, err
	}

	return validator, nil
}
