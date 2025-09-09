package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type (
	CustomValidator struct {
		validator  *validator.Validate
		translator ut.Translator
	}

	ValidationError struct {
		Field   string `json:"field"`
		Tag     string `json:"tag"`
		Value   string `json:"value"`
		Message string `json:"message"`
	}
)

// NewCustomValidator creates a new custom validator with English translations
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// Register custom tag name function
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Setup translator
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v, trans)

	return &CustomValidator{
		validator:  v,
		translator: trans,
	}
}

// Validate validates the struct and returns validation errors
// This implements Echo's Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		var errors []ValidationError

		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Value = fmt.Sprintf("%v", err.Value())
			element.Message = err.Translate(cv.translator)
			errors = append(errors, element)
		}

		return &ValidationErrors{Errors: errors}
	}
	return nil
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (ve *ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// GetValidationErrors returns the validation errors
func (ve *ValidationErrors) GetValidationErrors() []ValidationError {
	return ve.Errors
}
