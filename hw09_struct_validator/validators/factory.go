package validators

import "reflect"

func Factory(name string, structField reflect.StructField, constraint string) (Validator, error) {
	var validator Validator = nil
	switch name {
	case "max":
		validator = &MaxValidator{}
	case "min":
		validator = &MinValidator{}
	case "len":
		validator = &LenValidator{}
	case "in":
		validator = &InValidator{}
	default:
		return nil, ErrUnavailableValidator
	}

	if !isKindValid(validator, structField) {
		return nil, ErrInvalidType
	}

	err := validator.SetConstraint(constraint)
	if err != nil {
		return nil, err
	}

	return validator, nil
}

func isKindValid(validator Validator, structField reflect.StructField) bool {
	isKindValid := false
	validKinds := validator.GetValidKinds()
	for _, knd := range validKinds {
		rfType := structField.Type
		if rfType.Kind() == knd {
			isKindValid = true
			break
		}

		if rfType.Kind() == reflect.Slice && rfType.Elem().Kind() == knd {
			isKindValid = true
			break
		}
	}
	return isKindValid
}
