package validation

import (
	"fmt"
)

type BadRequest_FieldViolation struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

type BadRequest struct {
	FieldViolations []*BadRequest_FieldViolation
}

func fieldViolation(field string, err error) BadRequest_FieldViolation {
	return BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func validateStringNull(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("must not null value")
	}
	return nil
}

func validateInt32Null(value int32) error {
	if value == 0 {
		return fmt.Errorf("must not null value")
	}
	return nil
}

func validateString(value string, minLength, maxLength int) error {
	if err := validateStringNull(value); err != nil {
		return err
	}
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func validateInt(value int, minValue, maxValue int) error {
	if value < minValue || value > maxValue {
		return fmt.Errorf("value must between %d to %d", minValue, maxValue)
	}
	return nil
}

func validateInt32(value int32, minValue, maxValue int32) error {
	if value < minValue || value > maxValue {
		return fmt.Errorf("value must between %d to %d", minValue, maxValue)
	}
	return nil
}
