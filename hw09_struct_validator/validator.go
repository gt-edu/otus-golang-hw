package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gt-edu/otus-golang-hw/hw09_struct_validator/validators"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	strs := make([]string, len(v))
	for i, s := range v {
		strs[i] = fmt.Sprintf("%s: %s", s.Field, s.Err)
	}
	return "Validation errors: " + strings.Join(strs, "; ")
}

func Validate(v interface{}) error {
	rfVal := reflect.ValueOf(v)
	fieldsValidators, err := validators.ExtractValidators(rfVal)
	if err != nil {
		return err
	}

	var allErrors ValidationErrors
	for i := 0; i < rfVal.NumField(); i++ {
		validatorList := fieldsValidators[i]
		if validatorList == nil {
			continue
		}

		rfFieldVal := rfVal.Field(i)
		if !rfFieldVal.IsValid() || !rfFieldVal.CanInterface() {
			continue
		}

		rfFieldType := rfFieldVal.Type()
		for _, vld := range validatorList {
			if rfFieldType.Kind() == reflect.Slice {
				for ii := 0; ii < rfFieldVal.Len(); ii++ {
					err := vld.ValidateValue(rfFieldVal.Index(ii).Interface())
					if err != nil {
						allErrors = append(allErrors, ValidationError{Field: rfVal.Type().Field(i).Name, Err: err})
					}
				}
			} else {
				err := vld.ValidateValue(rfFieldVal.Interface())
				if err != nil {
					allErrors = append(allErrors, ValidationError{Field: rfVal.Type().Field(i).Name, Err: err})
				}
			}
		}
	}

	if len(allErrors) > 0 {
		return allErrors
	}
	return nil
}
