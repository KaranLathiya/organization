package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	error_handling "organization/error"
	"strings"

	"github.com/gookit/validate"
)

func BodyReadAndValidate(reader io.ReadCloser, bodyData interface{}, addValidationRules map[string]string) error {
	err := BodyRead(reader, bodyData)
	if err != nil {
		return err
	}
	err = ValidateStruct(bodyData, addValidationRules)
	if err != nil {
		return err
	}
	return nil
}

func BodyRead(reader io.ReadCloser, bodyData interface{}) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return error_handling.CreateCustomError(err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(body, &bodyData)
	if err != nil {
		return error_handling.UnmarshalError
	}
	return nil
}

func ValidateStruct(data interface{}, addValidationRules map[string]string) error {
	var errorMessage []string
	validator := validate.Struct(data)
	validator.StringRules(addValidationRules)
	if !(validator.Validate()) {
		var invalidDataArray []error_handling.InvalidData
		errors := validator.Errors.All()
		fmt.Println(errors) // all error messages
		for key, value := range errors {
			invalidData := error_handling.InvalidData{
				Field: key,
				Error: value,
			}
			invalidDataArray = append(invalidDataArray, invalidData)
			errorMessage = append(errorMessage, key)
		}
		return error_handling.CreateCustomError("Error in field:"+strings.Join(errorMessage, ","), http.StatusBadRequest, invalidDataArray...)
	}
	return nil
}
