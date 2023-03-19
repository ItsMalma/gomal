package gomal

import (
	"fmt"
	"net/mail"
	"reflect"
	"regexp"
	"unicode"
)

type Validator struct {
	name string

	value        any
	reflectValue reflect.Value
	valueType    reflect.Type

	errorMessages []string

	stop bool
}

func (validator Validator) getOption(option ...ValidatorOption) (ValidatorOption, bool) {
	if option == nil || len(option) < 1 {
		return ValidatorOption{}, false
	}
	return option[0], true
}

func (validator Validator) NotNil(option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	if validator.value == nil {
		opt, useOpt := validator.getOption(option...)
		if useOpt {
			validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
		} else {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must not be empty.", validator.name))
		}
	}
	return validator
}

func (validator Validator) NotEmpty(option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	errorMessage := ""
	if validator.value == nil {
		errorMessage = fmt.Sprintf("%v should not be empty.", validator.name)
		return validator
	} else {
		switch validator.valueType.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
			if validator.valueType.Kind() == reflect.Pointer && validator.reflectValue.Elem().Kind() != reflect.Array {
				return validator
			}
			if validator.reflectValue.Len() < 1 {
				errorMessage = fmt.Sprintf("%v should not be empty.", validator.name)
			}
		case reflect.Bool:
			if validator.value == false {
				errorMessage = fmt.Sprintf("%v should not be empty.", validator.name)
			}
		case reflect.Complex64, reflect.Complex128:
			if validator.value == 0+0i {
				errorMessage = fmt.Sprintf("%v should not be empty.", validator.name)
			}
		case reflect.Float32, reflect.Float64:
			if validator.value == 0.0 {
				errorMessage = fmt.Sprintf("%v should not be empty.", validator.name)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if validator.value == 0 {
				errorMessage = fmt.Sprintf("%v should not be empty.", validator.name)
			}
		case reflect.String:
			valueAsStr := validator.reflectValue.String()
			if len(valueAsStr) < 1 {
				errorMessage = fmt.Sprintf("%v should not be empty.", validator.name)
			} else {
				allWhitespace := true
				for _, ch := range valueAsStr {
					if !unicode.IsSpace(ch) {
						allWhitespace = false
					}
				}
				if allWhitespace {
					errorMessage = fmt.Sprintf("%v should not be empty.", validator.name)
				}
			}
		}
	}

	if errorMessage != "" {
		opt, useOpt := validator.getOption(option...)
		if useOpt {
			validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
		} else {
			validator.errorMessages = append(validator.errorMessages, errorMessage)
		}
	}

	return validator
}

func (validator Validator) NotEqual(another any, option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	if reflect.DeepEqual(validator.value, another) {
		opt, useOpt := validator.getOption(option...)
		if useOpt {
			validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
		} else {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should not be equal to %v.", validator.name, another))
		}
	}
	return validator
}

func (validator Validator) Equal(another any, option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	if !reflect.DeepEqual(validator.value, another) {
		opt, useOpt := validator.getOption(option...)
		if useOpt {
			validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
		} else {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should be equal to %v.", validator.name, another))
		}
	}
	return validator
}

// Only work for string
func (validator Validator) Length(min, max int, option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	if validator.valueType.Kind() == reflect.String {
		valueLength := validator.reflectValue.Len()
		if valueLength < min || valueLength > max {
			opt, useOpt := validator.getOption(option...)
			if useOpt {
				validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
			} else {
				validator.errorMessages = append(validator.errorMessages, fmt.Sprintf(
					"%v must be between %v and %v characters. You entered %v characters",
					validator.name, min, max, valueLength,
				))
			}
		}
	}
	return validator
}

// Only work for string
func (validator Validator) MaxLength(max int, option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	if validator.valueType.Kind() == reflect.String {
		valueLength := validator.reflectValue.Len()
		if valueLength > max {
			opt, useOpt := validator.getOption(option...)
			if useOpt {
				validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
			} else {
				validator.errorMessages = append(validator.errorMessages, fmt.Sprintf(
					"The length of %v must be %v characters or fewer. You entered %v characters.",
					validator.name, max, valueLength,
				))
			}
		}
	}
	return validator
}

// Only work for string
func (validator Validator) MinLength(min int, option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	if validator.valueType.Kind() == reflect.String {
		valueLength := validator.reflectValue.Len()
		if valueLength < min {
			opt, useOpt := validator.getOption(option...)
			if useOpt {
				validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
			} else {
				validator.errorMessages = append(validator.errorMessages, fmt.Sprintf(
					"The length of %v must be at least %v characters. You entered %v characters.",
					validator.name, min, valueLength,
				))
			}
		}
	}
	return validator
}

// Only work for numerical data type (int, uint, and float)
func (validator Validator) LessThan(another any, option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	errorMessage := ""
	switch validator.valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if validator.reflectValue.Int() >= another.(int64) {
			errorMessage = fmt.Sprintf("%v must be less than %v.", validator.name, another)
		}
	case reflect.Float32, reflect.Float64:
		if validator.reflectValue.Float() >= another.(float64) {
			errorMessage = fmt.Sprintf("%v must be less than %v.", validator.name, another)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if validator.reflectValue.Uint() >= another.(uint64) {
			errorMessage = fmt.Sprintf("%v must be less than %v.", validator.name, another)
		}
	}

	if errorMessage != "" {
		opt, useOpt := validator.getOption(option...)
		if useOpt {
			validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
		} else {
			validator.errorMessages = append(validator.errorMessages, errorMessage)
		}
	}

	return validator
}

// Only work for numerical data type (int, uint, and float)
func (validator Validator) LessThanOrEqual(another any, option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	errorMessage := ""
	switch validator.valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if validator.reflectValue.Int() > another.(int64) {
			errorMessage = fmt.Sprintf("%v must be less than or equal to %v.", validator.name, another)
		}
	case reflect.Float32, reflect.Float64:
		if validator.reflectValue.Float() > another.(float64) {
			errorMessage = fmt.Sprintf("%v must be less than or equal to %v.", validator.name, another)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if validator.reflectValue.Uint() > another.(uint64) {
			errorMessage = fmt.Sprintf("%v must be less than or equal to %v.", validator.name, another)
		}
	}

	if errorMessage != "" {
		opt, useOpt := validator.getOption(option...)
		if useOpt {
			validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
		} else {
			validator.errorMessages = append(validator.errorMessages, errorMessage)
		}
	}

	return validator
}

// Only work for numerical data type (int, uint, and float)
func (validator Validator) GreaterThan(another any, option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	errorMessage := ""
	switch validator.valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if validator.reflectValue.Int() <= another.(int64) {
			errorMessage = fmt.Sprintf("%v must be greater than %v.", validator.name, another)
		}
	case reflect.Float32, reflect.Float64:
		if validator.reflectValue.Float() <= another.(float64) {
			errorMessage = fmt.Sprintf("%v must be greater than %v.", validator.name, another)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if validator.reflectValue.Uint() <= another.(uint64) {
			errorMessage = fmt.Sprintf("%v must be greater than %v.", validator.name, another)
		}
	}

	if errorMessage != "" {
		opt, useOpt := validator.getOption(option...)
		if useOpt {
			validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
		} else {
			validator.errorMessages = append(validator.errorMessages, errorMessage)
		}
	}

	return validator
}

// Only work for numerical data type (int, uint, and float)
func (validator Validator) GreaterThanOrEqual(another any, option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	errorMessage := ""
	switch validator.valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if validator.reflectValue.Int() < another.(int64) {
			errorMessage = fmt.Sprintf("%v must be greater than or equal to %v.", validator.name, another)
		}
	case reflect.Float32, reflect.Float64:
		if validator.reflectValue.Float() < another.(float64) {
			errorMessage = fmt.Sprintf("%v must be greater than or equal to %v.", validator.name, another)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if validator.reflectValue.Uint() < another.(uint64) {
			errorMessage = fmt.Sprintf("%v must be greater than or equal to %v.", validator.name, another)
		}
	}

	if errorMessage != "" {
		opt, useOpt := validator.getOption(option...)
		if useOpt {
			validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
		} else {
			validator.errorMessages = append(validator.errorMessages, errorMessage)
		}
	}

	return validator
}

// Only work for string
func (validator Validator) RegExp(expr string, option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	if validator.valueType.Kind() == reflect.String {
		match, err := regexp.MatchString(expr, validator.reflectValue.String())
		if err != nil {
			panic(err)
		}
		if !match {
			opt, useOpt := validator.getOption(option...)
			if useOpt {
				validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
			} else {
				validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v is not in the correct format", validator.name))
			}
		}
	}

	return validator
}

// Only work for string
func (validator Validator) Email(option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	if validator.valueType.Kind() == reflect.String {
		if _, err := mail.ParseAddress(validator.reflectValue.String()); err != nil {
			opt, useOpt := validator.getOption(option...)
			if useOpt {
				validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
			} else {
				validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v is not a valid email address", validator.name))
			}
		}
	}

	return validator
}

func (validator Validator) Empty(option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	errorMessage := ""
	switch validator.valueType.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
		if validator.valueType.Kind() == reflect.Pointer && validator.reflectValue.Elem().Kind() != reflect.Array {
			return validator
		}
		if validator.reflectValue.Len() > 0 {
			errorMessage = fmt.Sprintf("%v must be empty", validator.name)
		}
	case reflect.Bool:
		if validator.value == false {
			errorMessage = fmt.Sprintf("%v must be empty", validator.name)
		}
	case reflect.Complex64, reflect.Complex128:
		if validator.value != 0+0i {
			errorMessage = fmt.Sprintf("%v must be empty", validator.name)
		}
	case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if validator.value != 0 {
			errorMessage = fmt.Sprintf("%v must be empty", validator.name)
		}
	case reflect.String:
		valueAsStr := validator.reflectValue.String()
		if len(valueAsStr) > 0 {
			errorMessage = fmt.Sprintf("%v must be empty", validator.name)
		}
		allWhitespace := true
		for _, ch := range valueAsStr {
			if !unicode.IsSpace(ch) {
				allWhitespace = false
			}
		}
		if !allWhitespace {
			errorMessage = fmt.Sprintf("%v must be empty", validator.name)
		}
	}

	if errorMessage != "" {
		opt, useOpt := validator.getOption(option...)
		if useOpt {
			validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
		} else {
			validator.errorMessages = append(validator.errorMessages, errorMessage)
		}
	}

	return validator
}

func (validator Validator) Nil(option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	if validator.value != nil {
		opt, useOpt := validator.getOption(option...)
		if useOpt {
			validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
		} else {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be empty.", validator.name))
		}
	}
	return validator
}

// Only work for numerical data type (int, uint, and float)
func (validator Validator) Between(min, max any, option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	errorMessage := ""
	switch validator.valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := validator.reflectValue.Int()
		minValue := min.(int64)
		maxValue := max.(int64)
		if value > minValue && value < maxValue {
			errorMessage = fmt.Sprintf("%v must be between %v and %v.", validator.name, minValue, maxValue)
		}
	case reflect.Float32, reflect.Float64:
		value := validator.reflectValue.Float()
		minValue := min.(float64)
		maxValue := max.(float64)
		if value > minValue && value < maxValue {
			errorMessage = fmt.Sprintf("%v must be between %v and %v.", validator.name, minValue, maxValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value := validator.reflectValue.Uint()
		minValue := min.(uint64)
		maxValue := max.(uint64)
		if value > minValue && value < maxValue {
			errorMessage = fmt.Sprintf("%v must be between %v and %v.", validator.name, minValue, maxValue)
		}
	}

	if errorMessage != "" {
		opt, useOpt := validator.getOption(option...)
		if useOpt {
			validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
		} else {
			validator.errorMessages = append(validator.errorMessages, errorMessage)
		}
	}

	return validator
}

// Unwrap if value is pointer
func (validator Validator) Unwrap() Validator {
	if validator.stop {
		return validator
	}

	if validator.reflectValue.Kind() == reflect.Pointer {
		validator.value = validator.reflectValue.Elem().Interface()
		validator.reflectValue = validator.reflectValue.Elem()
		validator.valueType = validator.reflectValue.Type()
	}

	return validator
}

func (validator Validator) Is(callback func() (bool, string), option ...ValidatorOption) Validator {
	if validator.stop {
		return validator
	}

	if success, errorMessage := callback(); !success {
		if errorMessage != "" {
			opt, useOpt := validator.getOption(option...)
			if useOpt {
				validator.errorMessages = append(validator.errorMessages, opt.ErrorMessage)
			} else {
				validator.errorMessages = append(validator.errorMessages, errorMessage)
			}
		}
	}

	return validator
}

func (validator Validator) When(condition bool) Validator {
	if !condition {
		validator.stop = true
	}
	return validator
}

func If(name string, value any) Validator {
	return Validator{
		name:          name,
		value:         value,
		reflectValue:  reflect.ValueOf(value),
		valueType:     reflect.TypeOf(value),
		errorMessages: []string{},
		stop:          false,
	}
}
