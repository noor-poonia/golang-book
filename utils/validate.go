package validate

import (
	// "fmt"
	// "go-rabbitmq/constants"
	// "errors"
	"fmt"
	"go-rabbitmq/model"

	"log"
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

type ISBNValidate struct {
	ISBN string `json:"isbn" validate:"isbn"`
}

type NameValidate struct {
	Name string `json:"name" validate:"alphaNumeric"`
}

type AuthorValidate struct {
	Author string `json:"author" validate:"alpha"`
}

func Validate(book model.Book) (string, error) {
	validate := validator.New()
	// str := ""
	
	validate.RegisterValidation("isbn", ValidateIsbn)
	validate.RegisterValidation("alphaNumeric", ValidateAlphaNumberic)
	validate.RegisterValidation("alpha", ValidateAlpha)
	isbn := ISBNValidate{
		ISBN: book.ISBN,
	}
	name := NameValidate{
		Name: book.Name,
	}
	author := AuthorValidate{
		Author: book.Author,
	}

	log.Printf("isbn value: %v\n", isbn)
	err := validate.Struct(isbn)
	// log.Printf("isbn struct error: %v\n", err.Error())
	if err != nil {
		return "ISBN is Invalid - Should be numeric with the length of 13", err
	}

	log.Printf("name value: %v\n", name)
	er := validate.Struct(name)
	// log.Printf("name struct error: %v\n", er.Error())
	if er != nil {
		return "Book Name is Invalid - Should not contain anything other than alphabets and Numbers", er
	}

	log.Printf("author value: %v\n", author)
	e := validate.Struct(author)
	// log.Printf("author struct error: %v\n", e.Error())
	if e != nil {
		return "Author Name is Invalid - Should not contain anything other than alphabets", e
	}
	
	// e1 := validate.StructExcept(book, book.Author, book.ISBN)
	// log.Printf("value of e1: %v\n", e1.Error())
	// if e1 != nil {
	// 	str = "Book Name is Invalid - Should not contain anything other than alphabets and Numbers"
	// 	fmt.Printf("returning from validation: %s\n", str)
	// 	return str
	// } 
	
	// e2 := validate.StructExcept(book, book.Name, book.Author)
	// log.Printf("value of e2: %v\n", e2.Error())
	// if e2 != nil {
	// 	str = "ISBN is Invalid - Should be numeric with the length of 13"
	// 	fmt.Printf("returning from validation: %s\n", str)
	// 	return str
	// } 
	
	// e3 := validate.StructExcept(book, book.Name, book.ISBN)
	// log.Printf("value of e3: %v\n", e3.Error())
	// if e3 != nil {
	// 	str = "Author Name is Invalid - Should not contain anything other than alphabets"
	// 	fmt.Printf("returning from validation: %s\n", str)
	// 	// return str, e
	// 	return str
	// } 
	return "", nil
}

func ValidateAlpha(f1 validator.FieldLevel) (bool) {
	text := f1.Field().String()
	alphaRegex := regexp.MustCompile(`^[a-zA-Z\s.]+$`)
	return alphaRegex.MatchString(text)
}

func ValidateAlphaNumberic(f1 validator.FieldLevel) (bool) {
	text := f1.Field().String()
	alphaNumRegex := regexp.MustCompile(`^[A-Za-z0-9\s.']+$`)
	return alphaNumRegex.MatchString(text)
}

func ValidateIsbn(f1 validator.FieldLevel) (bool) {
	text := f1.Field().String()
	fmt.Printf("validateISBN text: %v\n", text)
	if len(text) == 13 {
		val := regexp.MustCompile(`\d`).MatchString(text)
		fmt.Printf("validateISBN result: %v\n", val)
		return val
	}
	return false
}