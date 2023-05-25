package model

type Book struct {
	Name string `json:"name" validate:"required,min=3,max=50,alphaNumeric"`
	Author string `json:"author" validate:"required,min=3,max=20,alpha"`
	Pages int `json:"pages" validate:"required,numeric"`
}
