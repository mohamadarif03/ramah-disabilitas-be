package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors { // Iterate over all errors
			field := e.Field()
			// Basic Name Transformation (CamelCase -> snake_case somewhat)
			fieldName := strings.ToLower(field)

			// Custom overrides based on common fields
			switch field {
			case "ConfirmPassword":
				fieldName = "confirm_password"
				field = "Konfirmasi Password"
			case "Password":
				fieldName = "password"
				field = "Password"
			case "Email":
				fieldName = "email"
				field = "Email"
			case "Name":
				fieldName = "name"
				field = "Nama"
			case "Role":
				fieldName = "role"
				field = "Peran"
			case "Title":
				fieldName = "title"
				field = "Judul"
			case "Description":
				fieldName = "description"
				field = "Deskripsi"
			case "Thumbnail":
				fieldName = "thumbnail"
				field = "Thumbnail"
			case "ClassCode":
				fieldName = "class_code"
				field = "Kode Kelas"
			}

			// Custom message logic
			var msg string
			switch e.Tag() {
			case "required":
				msg = fmt.Sprintf("%s wajib diisi.", field)
			case "email":
				msg = fmt.Sprintf("Format %s tidak valid.", field)
			case "min":
				msg = fmt.Sprintf("%s minimal %s karakter.", field, e.Param())
			case "max":
				msg = fmt.Sprintf("%s maksimal %s karakter.", field, e.Param())
			case "eqfield":
				msg = fmt.Sprintf("%s harus sama dengan %s.", field, e.Param())
			default:
				msg = fmt.Sprintf("%s tidak valid.", field)
			}

			errors[fieldName] = msg
		}
		return errors
	}

	// Fallback if it's not a validation error
	errors["general"] = err.Error()
	return errors
}
