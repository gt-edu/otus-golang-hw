package validators

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func ExtractValidators(rfVal reflect.Value) ([][]Validator, error) {
	if rfVal.Kind() != reflect.Struct {
		return nil, ErrValueIsNotStruct
	}

	rfType := rfVal.Type()
	fieldsValidators := make([][]Validator, rfType.NumField())
	for i := 0; i < rfType.NumField(); i++ {
		validateTag, ok := rfType.Field(i).Tag.Lookup("validate")
		if !ok {
			fieldsValidators[i] = nil
			continue
		}

		validatorDescriptors := strings.Split(validateTag, "|")
		var oneFieldValidators []Validator
		for _, validatorDescriptor := range validatorDescriptors {
			validatorData := strings.Split(validatorDescriptor, ":")
			if len(validatorData) != 2 || len(validatorData[0]) == 0 || len(validatorData[1]) == 0 {
				return nil, errors.Wrap(ErrInvalidValidateTag, fmt.Sprintf("invalid value '%s'", validatorDescriptor))
			}

			validator, err := Factory(validatorData[0], rfType.Field(i), validatorData[1])
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("'%s'", validatorData[0]))
			}

			oneFieldValidators = append(oneFieldValidators, validator)
		}

		fieldsValidators[i] = oneFieldValidators
	}
	return fieldsValidators, nil
}
