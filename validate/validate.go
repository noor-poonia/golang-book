package validate

import (
	"fmt"
	"go-rabbitmq/model"
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

func Validate(book model.Book) error {
	validate := validator.New()

	err := validate.RegisterValidation("alphaNumeric", ValidateAlphaNumberic)
	if err != nil {
		fmt.Printf("err msg: %v\n", err)
		// error message
	}
	e := validate.RegisterValidation("alpha", ValidateAlpha)
	if e != nil {
		fmt.Printf("err msg: %v\n", e)
        // error message
    }

	return validate.Struct(book)
}

func ValidateAlpha(f1 validator.FieldLevel) bool {
	text := f1.Field().String()
	alphaRegex := regexp.MustCompile(`^[a-zA-Z\s.]+$`)
	return alphaRegex.MatchString(text)
}

func ValidateAlphaNumberic(f1 validator.FieldLevel) bool {
	text := f1.Field().String()
	alphaNumRegex := regexp.MustCompile(`^[A-Za-z0-9\s.']+$`)
	return alphaNumRegex.MatchString(text)
}