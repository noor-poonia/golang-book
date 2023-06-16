package constants

var Logs = map[string]string{
	"A100": "Data Successfully Published",
	"A101": "Data Received Successfully",
	"A102" : "No Such Book",
	"A104" : "Failed to Retrieve Books",
	"A105": "Failed to Parse Data",
	"A106": "Data Validation Failed",
	"A107": "Successfully Deleted Book",
	"A108": "ISBN Number Not Found",
	"A109" : "Book Updated",
	"A110" : "Failed to Update Book",
	"A111": "Author Name is Invalid - Should not contain anything other than alphabets",
	"A112": "Book Name is Invalid - Should not contain anything other than alphabets and Numbers",
	"A113": "ISBN is Invalid - Should be numeric with the length of 13",
	"A114": "Check RabbitMQ Connection",
	"A115" : "Book Already Exists",
}

// add validation msg