package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Person struct {
	ID          string `json:"id" binding:"required,number"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phoneNumber" binding:"required,e164,len=12"`
	CountryCode string `json:"countryCode" binding:"required,iso3166_1_alpha2"`
	SSN         string `json:"ssn" binding:"required,ssn,len=11"`
}

// FieldErrorMsg: custom error message for a FieldError post validation
type FieldErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func main() {

	router := gin.Default()

	router.POST("/register", register)

	router.Run(":8000")
}

func register(ctx *gin.Context) {

	// expected request body input
	var person Person

	// check for errors and proccess validations then read request body into struct
	if err := ctx.ShouldBindJSON(&person); err != nil {

		// declare ValidationErrors variable
		var validationErrs validator.ValidationErrors

		// check if received error equals ValidationErrors | if so initialize validationErrs variable with ValidationErrors Slice
		if errors.As(err, &validationErrs) {

			// create slice the length of ValidationErrs
			out := make([]FieldErrorMsg, len(validationErrs))

			// for each field that failed validation create a struct with custom error
			for i, fieldErr := range validationErrs {
				out[i] = FieldErrorMsg{
					Field:   fieldErr.Field(),
					Message: getErrorMsg(fieldErr),
				}
			}

			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": out})

			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	fmt.Printf("Person Data:%+v\n\n", person)

	ctx.JSON(http.StatusOK, gin.H{"person": &person})
}

// getErrorMsg: returns an error message for a failed validation tag associated to a FieldError
func getErrorMsg(fieldErr validator.FieldError) string {

	// return error based on the validation tag that failed
	switch failedTag := fieldErr.Tag(); failedTag {

	case "required":
		return "this field is required"

	case "e164":
		return fmt.Sprintln("expected e164 phone number format e.g. +11234567890")

	case "email":
		return fmt.Sprintf("invalid email format expected format baki@gmail.com received: %v", fieldErr.Value())

	case "iso3166_1_alpha2":
		return fmt.Sprintf("invalid country code format expected two capital letters received: %v", fieldErr.Value())

	case "min":
		return fmt.Sprintf("field minimum length: %v received field length: %v", fieldErr.Param(), len(fieldErr.Value().(string)))

	case "max":
		return fmt.Sprintf("field maximum length: %v received field length: %v", fieldErr.Param(), len(fieldErr.Value().(string)))

	case "len":
		return fmt.Sprintf("required field length: %v received field length: %v", fieldErr.Param(), len(fieldErr.Value().(string)))

	case "number":
		return fmt.Sprintf("expected numbers 0-9 received: %v", fieldErr.Value())

	case "ssn":
		return fmt.Sprintf("invalid ssn format expected: 123-456-6789")

	default:
		return "unkown error"
	}
}

// gin Binding

//   - a deserialization library part of the gin web framework with built-in struct validation
//   - uses the validator package and struct tags for value validation of structs and fields
//   - deserializes an http request body reading the body buffer into structs and maps
