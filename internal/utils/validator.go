package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

func ValidateStruct(s interface{}) ValidationErrors {
	var validationErrors ValidationErrors

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   fmt.Sprintf("%v", err.Value()),
				Message: getErrorMessage(err),
			})
		}
	}

	return validationErrors
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("Поле %s обязательно для заполнения", fe.Field())
	case "min":
		return fmt.Sprintf("Поле %s должно содержать минимум %s символов", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("Поле %s должно содержать максимум %s символов", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("Поле %s должно содержать корректный email адрес", fe.Field())
	case "alphanum":
		return fmt.Sprintf("Поле %s должно содержать только буквы и цифры", fe.Field())
	default:
		return fmt.Sprintf("Поле %s некорректно", fe.Field())
	}
}
