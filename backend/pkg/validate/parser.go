package validate

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Parse converts validator.ValidationErrors into a slice of FieldError.
// Returns nil if err is not a ValidationErrors type.
func Parse(err error) *FieldError {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}

	// out := make([]FieldError, len(ve))
	// for i, fe := range ve {
	// 	out[i] = FieldError{
	// 		Field:   fe.Field(),
	// 		Message: msgForTag(fe),
	// 	}
	// }

	// Only take the first error
	fe := ve[0]
	return &FieldError{
		Field:   fe.Field(),
		Message: msgForTag(fe),
	}
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", fe.Field())
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s harus lebih besar atau sama dengan %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s harus lebih kecil atau sama dengan %s", fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf("%s harus lebih besar dari %s", fe.Field(), fe.Param())
	case "lt":
		return fmt.Sprintf("%s harus lebih kecil dari %s", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("%s harus berupa alamat email yang valid", fe.Field())
	case "alphanum":
		return fmt.Sprintf("%s hanya boleh mengandung huruf dan angka", fe.Field())
	case "numeric":
		return fmt.Sprintf("%s harus berupa angka", fe.Field())
	case "oneof":
		return fmt.Sprintf("%s harus salah satu dari: %s", fe.Field(), fe.Param())
	case "url":
		return fmt.Sprintf("%s harus berupa URL yang valid", fe.Field())
	default:
		return fmt.Sprintf("%s tidak valid", fe.Field())
	}
}
