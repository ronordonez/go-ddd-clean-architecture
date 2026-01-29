package application

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"go-architecture/internal/shared/errors"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		details := make(map[string]interface{})
		
		for _, fieldErr := range validationErrors {
			details[fieldErr.Field()] = fmt.Sprintf(
				"Field validation for '%s' failed on the '%s' tag",
				fieldErr.Field(),
				fieldErr.Tag(),
			)
		}
		
		return errors.NewValidationError("Validation failed", details)
	}
	return nil
}
