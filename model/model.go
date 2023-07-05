package model

type Book struct {
	Name string `json:"name" validate:"required,min=3,max=50,alphaNumeric"`
	ISBN string `json:"isbn" validate:"required"`
	Author string `json:"author" validate:"required,min=3,max=20,alpha"`
	Pages int `json:"pages" validate:"required,numeric"`
}

type Response struct {
	Code string `json:"code"`
	Message string `json:"message"`
	Data []Book `json:"data"`
	Error []string `json:"error"`
}