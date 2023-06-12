package validate

import (
	"go-rabbitmq/model"
	"log"
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

func Validate(book model.Book) error {
	validate := validator.New()

	err := validate.RegisterValidation("alphaNumeric", ValidateAlphaNumberic)
	if err != nil {
		log.Printf("err msg: %v\n", err)
		// error message
	}
	e := validate.RegisterValidation("alpha", ValidateAlpha)
	if e != nil {
		log.Printf("err msg: %v\n", e)
        // error message
    }
	er := validate.RegisterValidation("isbn", ValidateIsbn)
	if er != nil {
		log.Printf("err msg: %v\n", e)
        // error message
    }
	return validate.Struct(book)
}

func ValidateAlpha(f1 validator.FieldLevel) (bool) {
	text := f1.Field().String()
	alphaRegex := regexp.MustCompile(`^[a-zA-Z\s.]+$`)
	return alphaRegex.MatchString(text)
	// if !res {
	// 	return false, 
	// }
	// return error message also
}

func ValidateAlphaNumberic(f1 validator.FieldLevel) bool {
	text := f1.Field().String()
	alphaNumRegex := regexp.MustCompile(`^[A-Za-z0-9\s.']+$`)
	return alphaNumRegex.MatchString(text)
}

func ValidateIsbn(f1 validator.FieldLevel) bool {
	text := f1.Field().String()
	if len(text) == 13 {
		return regexp.MustCompile(`\d`).MatchString(text)
	}
	return false
}

// func ValidatePages(f1 validator.FieldLevel) bool {
// 	text := f1.Field().String()
// 	return regexp.MustCompile(`\d`).MatchString(text)
// }