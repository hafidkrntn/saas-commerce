package validate

import (
	"backend-go/pkg/apperror"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitValidator registers the JSON tag name func into Gin's built-in validator.
// Call this once in main.go before starting the server.
func InitValidator() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	// Use json tag name (e.g. "warehouse_id") instead of struct field name ("WarehouseID")
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate validates the request body on the first occurrence.
func Validate(c *gin.Context, body any) error {
	if err := c.ShouldBindJSON(body); err != nil {
		if field := Parse(err); field != nil {
			return apperror.BadRequest(field.Message, err)
		}
		return apperror.BadRequest("request body tidak valid", err)
	}
	return nil
}

func ParseBindError(err error) error {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return apperror.BadRequest("request body tidak valid", err)
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return apperror.BadRequest(fmt.Sprintf("kolom '%s' harus merupakan %s, bukan %s", typeErr.Field, typeErr.Type, typeErr.Value), err)
	}

	msg := err.Error()
	if strings.Contains(msg, "UUID") || strings.Contains(msg, "uuid") {
		return apperror.BadRequest(fmt.Sprintf("uuid tidak valid: %s", msg), err)
	}

	// 4. Validator errors (field-level validation tags)
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		var messages []string
		for _, e := range validationErrs {
			messages = append(messages, fmt.Sprintf(
				"kolom '%s' gagal untuk validasi '%s'", e.Field(), e.Tag(),
			))
		}
		return apperror.BadRequest(strings.Join(messages, "; "), err)
	}

	// 5. Fallback
	return err
}
