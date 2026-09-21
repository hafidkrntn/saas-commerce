package apperror

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime"
)

// AppError is a structured error that carries an HTTP status code and a user-facing message.
type AppError struct {
	err        error
	message    string
	statusCode int
}

// Field is a key-value pair for adding context to error logs.
type Field struct {
	Key   string
	Value any
}

// --- Constructors ---

// New creates an AppError and logs it with caller info.
func New(statusCode int, message string, err error, fields ...Field) *AppError {
	pc, file, line, _ := runtime.Caller(1)
	trace := fmt.Sprintf("%s:%d", path.Base(file), line)
	caller := path.Base(runtime.FuncForPC(pc).Name())

	if err != nil {
		log.Printf("[ERROR] %s %s | %s | %v | fields=%v", trace, caller, message, err, fields)
	} else {
		log.Printf("[ERROR] %s %s | %s | fields=%v", trace, caller, message, fields)
	}

	return &AppError{
		statusCode: statusCode,
		message:    message,
		err:        errors.New(message),
	}
}

func BadRequest(message string, err error, fields ...Field) *AppError {
	return New(http.StatusBadRequest, message, err, fields...)
}

func NotFound(message string, err error, fields ...Field) *AppError {
	return New(http.StatusNotFound, message, err, fields...)
}

func Unauthorized(message string, err error, fields ...Field) *AppError {
	return New(http.StatusUnauthorized, message, err, fields...)
}

func Forbidden(message string, err error, fields ...Field) *AppError {
	return New(http.StatusForbidden, message, err, fields...)
}

func Internal(message string, err error, fields ...Field) *AppError {
	return New(http.StatusInternalServerError, message, err, fields...)
}

// --- Interface implementations ---

func (e *AppError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return e.message
}

func (e *AppError) Message() string {
	return e.message
}

func (e *AppError) StatusCode() int {
	return e.statusCode
}

func (e *AppError) Unwrap() error {
	return e.err
}

// --- Helper to create Field ---

func F(key string, value any) Field {
	return Field{Key: key, Value: value}
}
