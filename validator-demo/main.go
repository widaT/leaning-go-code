package main

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	Name      	   string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=150"`
	Email          string     `validate:"required,email"`
}

var validate *validator.Validate

func main() {

	validate = validator.New()

	validateStruct()
	//validateVariable()
}

func validateStruct() {
	user := &User{
		Name:      		"wida",
		Age:            165,
		Email:          "someone.gmail.com",
	}

	err := validate.Struct(user)
	fmt.Println(err)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace()) // can differ when a custom TagNameFunc is registered or
			fmt.Println(err.StructField())     // by passing alt name to ReportError like below
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
		return
	}
}

func validateVariable() {
	myEmail := "wida59.gmail.com"
	errs := validate.Var(myEmail, "required,email")
	if errs != nil {
		fmt.Println(errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
		return
	}
}