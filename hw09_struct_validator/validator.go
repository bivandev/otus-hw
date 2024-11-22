package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, e := range v {
		sb.WriteString(fmt.Sprintf("field '%s': %v\n", e.Field, e.Err))
	}

	return sb.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %s", val.Kind())
	}

	var errs ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		fieldErrors := validateField(field.Name, fieldValue, validateTag)
		if len(fieldErrors) > 0 {
			errs = append(errs, fieldErrors...)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateField(fieldName string, value reflect.Value, tag string) ValidationErrors {
	var errs ValidationErrors
	validators := strings.Split(tag, "|")

	for _, validator := range validators {
		parts := strings.SplitN(validator, ":", 2)
		if len(parts) != 2 {
			errs = append(errs, ValidationError{Field: fieldName, Err: fmt.Errorf("invalid validator format: %s", validator)})
			continue
		}

		var (
			err   error
			param = parts[1]
			vName = parts[0]
		)

		switch value.Kind() { //nolint:exhaustive
		case reflect.String:
			err = validateString(value.String(), vName, param)
		case reflect.Int:
			err = validateInt(int(value.Int()), vName, param)
		case reflect.Slice:
			err = validateSlice(value, vName, param)
		default:
			err = fmt.Errorf("unsupported type: %s", value.Kind())
		}

		if err != nil {
			errs = append(errs, ValidationError{Field: fieldName, Err: err})
		}
	}

	return errs
}

func validateString(value, validator, param string) error {
	switch validator {
	case "len":
		length, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("invalid length parameter: %w", err)
		}

		if len(value) != length {
			return fmt.Errorf("must be %d characters long", length)
		}
	case "regexp":
		re, err := regexp.Compile(param)
		if err != nil {
			return fmt.Errorf("invalid regexp: %w", err)
		}

		if !re.MatchString(value) {
			return fmt.Errorf("does not match pattern %s", param)
		}
	case "in":
		allowedValues := strings.Split(param, ",")
		for _, v := range allowedValues {
			if value == v {
				return nil
			}
		}

		return fmt.Errorf("must be one of %s", strings.Join(allowedValues, ", "))
	default:
		return fmt.Errorf("unknown string validator: %s", validator)
	}

	return nil
}

func validateInt(value int, validator, param string) error {
	switch validator {
	case "min":
		mn, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("invalid min parameter: %w", err)
		}

		if value < mn {
			return fmt.Errorf("must be >= %d", mn)
		}
	case "max":
		mx, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("invalid max parameter: %w", err)
		}

		if value > mx {
			return fmt.Errorf("must be <= %d", mx)
		}
	case "in":
		allowedValues := strings.Split(param, ",")
		for _, v := range allowedValues {
			if iv, err := strconv.Atoi(v); err == nil && value == iv {
				return nil
			}
		}

		return fmt.Errorf("must be one of %s", strings.Join(allowedValues, ", "))
	default:
		return fmt.Errorf("unknown int validator: %s", validator)
	}

	return nil
}

func validateSlice(value reflect.Value, validator, param string) error {
	for i := 0; i < value.Len(); i++ {
		var (
			elem = value.Index(i)
			err  error
		)

		switch elem.Kind() { //nolint:exhaustive
		case reflect.String:
			err = validateString(elem.String(), validator, param)
		case reflect.Int:
			err = validateInt(int(elem.Int()), validator, param)
		default:
			return fmt.Errorf("unsupported slice element type: %s", elem.Kind())
		}

		if err != nil {
			return fmt.Errorf("element %d: %w", i, err)
		}
	}

	return nil
}
