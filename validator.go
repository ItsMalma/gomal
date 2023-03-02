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

func (validator Validator) NotNil() Validator {
	if validator.stop {
		return validator
	}

	if validator.value == nil {
		validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must not be empty.", validator.name))
	}
	return validator
}

func (validator Validator) NotEmpty() Validator {
	if validator.stop {
		return validator
	}

	if validator.value == nil {
		validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should not be empty.", validator.name))
		return validator
	} else {
		switch validator.valueType.Kind() {
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
			if validator.valueType.Kind() == reflect.Pointer && validator.reflectValue.Elem().Kind() != reflect.Array {
				return validator
			}
			if validator.reflectValue.Len() < 1 {
				validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should not be empty.", validator.name))
			}
		case reflect.Bool:
			if validator.value == false {
				validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should not be empty.", validator.name))
			}
		case reflect.Complex64, reflect.Complex128:
			if validator.value == 0+0i {
				validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should not be empty.", validator.name))
			}
		case reflect.Float32, reflect.Float64:
			if validator.value == 0.0 {
				validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should not be empty.", validator.name))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if validator.value == 0 {
				validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should not be empty.", validator.name))
			}
		case reflect.String:
			valueAsStr := validator.reflectValue.String()
			if len(valueAsStr) < 1 {
				validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should not be empty.", validator.name))
			} else {
				allWhitespace := true
				for _, ch := range valueAsStr {
					if !unicode.IsSpace(ch) {
						allWhitespace = false
					}
				}
				if allWhitespace {
					validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should not be empty.", validator.name))
				}
			}
		}
	}

	return validator
}

func (validator Validator) NotEqual(another any) Validator {
	if validator.stop {
		return validator
	}

	if reflect.DeepEqual(validator.value, another) {
		validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should not be equal to %v.", validator.name, another))
	}
	return validator
}

func (validator Validator) Equal(another any) Validator {
	if validator.stop {
		return validator
	}

	if !reflect.DeepEqual(validator.value, another) {
		validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v should be equal to %v.", validator.name, another))
	}
	return validator
}

// Only work for string
func (validator Validator) Length(min, max int) Validator {
	if validator.stop {
		return validator
	}

	if validator.valueType.Kind() == reflect.String {
		valueLength := validator.reflectValue.Len()
		if valueLength < min || valueLength > max {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf(
				"%v must be between %v and %v characters. You entered %v characters",
				validator.name, min, max, valueLength,
			))
		}
	}
	return validator
}

// Only work for string
func (validator Validator) MaxLength(max int) Validator {
	if validator.stop {
		return validator
	}

	if validator.valueType.Kind() == reflect.String {
		valueLength := validator.reflectValue.Len()
		if valueLength > max {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf(
				"The length of %v must be %v characters or fewer. You entered %v characters.",
				validator.name, max, valueLength,
			))
		}
	}
	return validator
}

// Only work for string
func (validator Validator) MinLength(min int) Validator {
	if validator.stop {
		return validator
	}

	if validator.valueType.Kind() == reflect.String {
		valueLength := validator.reflectValue.Len()
		if valueLength < min {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf(
				"The length of %v must be at least %v characters. You entered %v characters.",
				validator.name, min, valueLength,
			))
		}
	}
	return validator
}

// Only work for numerical data type (int, uint, and float)
func (validator Validator) LessThan(another any) Validator {
	if validator.stop {
		return validator
	}

	switch validator.valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if validator.reflectValue.Int() >= another.(int64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be less than %v.", validator.name, another))
		}
	case reflect.Float32, reflect.Float64:
		if validator.reflectValue.Float() >= another.(float64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be less than %v.", validator.name, another))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if validator.reflectValue.Uint() >= another.(uint64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be less than %v.", validator.name, another))
		}
	}
	return validator
}

// Only work for numerical data type (int, uint, and float)
func (validator Validator) LessThanOrEqual(another any) Validator {
	if validator.stop {
		return validator
	}

	switch validator.valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if validator.reflectValue.Int() > another.(int64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be less than or equal to %v.", validator.name, another))
		}
	case reflect.Float32, reflect.Float64:
		if validator.reflectValue.Float() > another.(float64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be less than or equal to %v.", validator.name, another))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if validator.reflectValue.Uint() > another.(uint64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be less than or equal to %v.", validator.name, another))
		}
	}
	return validator
}

// Only work for numerical data type (int, uint, and float)
func (validator Validator) GreaterThan(another any) Validator {
	if validator.stop {
		return validator
	}

	switch validator.valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if validator.reflectValue.Int() <= another.(int64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be greater than %v.", validator.name, another))
		}
	case reflect.Float32, reflect.Float64:
		if validator.reflectValue.Float() <= another.(float64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be greater than %v.", validator.name, another))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if validator.reflectValue.Uint() <= another.(uint64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be greater than %v.", validator.name, another))
		}
	}
	return validator
}

// Only work for numerical data type (int, uint, and float)
func (validator Validator) GreaterThanOrEqual(another any) Validator {
	if validator.stop {
		return validator
	}

	switch validator.valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if validator.reflectValue.Int() < another.(int64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be greater than or equal to %v.", validator.name, another))
		}
	case reflect.Float32, reflect.Float64:
		if validator.reflectValue.Float() < another.(float64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be greater than or equal to %v.", validator.name, another))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if validator.reflectValue.Uint() < another.(uint64) {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be greater than or equal to %v.", validator.name, another))
		}
	}
	return validator
}

// Only work for string
func (validator Validator) RegExp(expr string) Validator {
	if validator.stop {
		return validator
	}

	if validator.valueType.Kind() == reflect.String {
		match, err := regexp.MatchString(expr, validator.reflectValue.String())
		if err != nil {
			panic(err)
		}
		if !match {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v is not in the correct format", validator.name))
		}
	}

	return validator
}

// Only work for string
func (validator Validator) Email() Validator {
	if validator.stop {
		return validator
	}

	if validator.valueType.Kind() == reflect.String {
		if _, err := mail.ParseAddress(validator.reflectValue.String()); err != nil {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v is not a valid email address", validator.name))
		}
	}

	return validator
}

func (validator Validator) Empty() Validator {
	if validator.stop {
		return validator
	}

	switch validator.valueType.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Pointer, reflect.Slice:
		if validator.valueType.Kind() == reflect.Pointer && validator.reflectValue.Elem().Kind() != reflect.Array {
			return validator
		}
		if validator.reflectValue.Len() > 0 {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be empty", validator.name))
		}
	case reflect.Bool:
		if validator.value == false {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be empty", validator.name))
		}
	case reflect.Complex64, reflect.Complex128:
		if validator.value != 0+0i {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be empty", validator.name))
		}
	case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if validator.value != 0 {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be empty", validator.name))
		}
	case reflect.String:
		valueAsStr := validator.reflectValue.String()
		if len(valueAsStr) > 0 {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be empty", validator.name))
		}
		allWhitespace := true
		for _, ch := range valueAsStr {
			if !unicode.IsSpace(ch) {
				allWhitespace = false
			}
		}
		if !allWhitespace {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be empty", validator.name))
		}
	}
	return validator
}

func (validator Validator) Nil() Validator {
	if validator.stop {
		return validator
	}

	if validator.value != nil {
		validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be empty.", validator.name))
	}
	return validator
}

// Only work for numerical data type (int, uint, and float)
func (validator Validator) Between(min, max any) Validator {
	if validator.stop {
		return validator
	}

	switch validator.valueType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := validator.reflectValue.Int()
		minValue := min.(int64)
		maxValue := max.(int64)
		if value > minValue && value < maxValue {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be between %v and %v.", validator.name, minValue, maxValue))
		}
	case reflect.Float32, reflect.Float64:
		value := validator.reflectValue.Float()
		minValue := min.(float64)
		maxValue := max.(float64)
		if value > minValue && value < maxValue {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be between %v and %v.", validator.name, minValue, maxValue))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value := validator.reflectValue.Uint()
		minValue := min.(uint64)
		maxValue := max.(uint64)
		if value > minValue && value < maxValue {
			validator.errorMessages = append(validator.errorMessages, fmt.Sprintf("%v must be between %v and %v.", validator.name, minValue, maxValue))
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
