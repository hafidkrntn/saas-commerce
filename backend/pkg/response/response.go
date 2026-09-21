package response

import (
	"backend-go/constants"
	"backend-go/pkg/apperror"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ==================== GLOBAL DB ====================

var DB *gorm.DB

// SetDB initializes the DB instance for activity logging.
func SetDB(db *gorm.DB) {
	DB = db
}

// ==================== STRUCT ====================

type Response struct {
	Code  int    `json:"code"`
	Data  any    `json:"data,omitempty"`
	Msg   string `json:"msg"`
	Error string `json:"error,omitempty"`
}

// ==================== CORE ====================

func JSON(c *gin.Context, code int, data any, msg string) {
	c.JSON(code, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
	c.Abort()
}

// ==================== ERROR ====================

// Err handles any error and sends the appropriate HTTP response.
func Err(c *gin.Context, err error) {
	if err == nil {
		JSON(c, http.StatusInternalServerError, nil, constants.ErrInternalServer)
		return
	}

	// JSON parse/unmarshal errors → 400
	var syn *json.SyntaxError
	var unm *json.UnmarshalTypeError
	if errors.As(err, &syn) || errors.As(err, &unm) {
		JSON(c, http.StatusBadRequest, nil, constants.ErrInvalidInput)
		return
	}

	// AppError → use its status code and message
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		JSON(c, appErr.StatusCode(), nil, appErr.Message())
		return
	}

	// Validation errors from gin binding
	if handleValidationError(c, err) {
		return
	}

	// GORM not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		JSON(c, http.StatusNotFound, nil, constants.ErrNotFound)
		return
	}

	// Fallback
	JSON(c, http.StatusInternalServerError, nil, constants.ErrInternalServer)
}

func ErrBadRequest(c *gin.Context, err error) {
	if err == nil {
		return
	}
	JSON(c, http.StatusBadRequest, nil, constants.ErrInvalidInput)
}

func ErrUnauthorized(c *gin.Context, err error) {
	if err == nil {
		return
	}
	JSON(c, http.StatusUnauthorized, nil, constants.ErrTokenInvalid)
}

func ErrForbidden(c *gin.Context, err error) {
	if err == nil {
		return
	}
	JSON(c, http.StatusForbidden, nil, constants.ErrForbidden)
}

// ==================== SUCCESS ====================

func OK(c *gin.Context, data any) {
	JSON(c, http.StatusOK, data, constants.MsgDataFound)
}

func Created(c *gin.Context, data any, payload any) {
	logActivity(c, payload)
	JSON(c, http.StatusCreated, data, constants.MsgDataAdd)
}

func Updated(c *gin.Context, data any, payload any) {
	logActivity(c, payload)
	JSON(c, http.StatusOK, data, constants.MsgDataUpdate)
}

func Deleted(c *gin.Context, data any) {
	logActivity(c, data)
	JSON(c, http.StatusOK, nil, constants.MsgDataDelete)
}

// ==================== HELPER ====================

func WrapErr(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

func logActivity(c *gin.Context, data any) {
	// Get module name from caller
	module := ""
	pc, _, _, ok := runtime.Caller(2)
	if ok {
		fn := runtime.FuncForPC(pc).Name()
		lastSlash := strings.LastIndexByte(fn, '/')
		fn = fn[lastSlash+1:]
		dot := strings.IndexByte(fn, '.')
		if dot > 0 {
			module = fn[:dot]
		}
	}

	method := c.Request.Method
	if method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch && method != http.MethodDelete {
		return
	}

	if DB == nil {
		return
	}

	// Capture values from gin.Context BEFORE goroutine (context may be recycled)
	userIDStr := c.GetString("UserId")
	var userID *string
	if userIDStr != "" {
		if _, err := uuid.Parse(userIDStr); err == nil {
			copied := userIDStr
			userID = &copied
		}
	}

	recordIDParam := c.Param("id")
	var recordID *string
	if recordIDParam != "" {
		if _, err := uuid.Parse(recordIDParam); err == nil {
			copied := recordIDParam
			recordID = &copied
		}
	}

	var payloadStr *string
	if data != nil {
		if b, err := json.Marshal(data); err == nil {
			s := string(b)
			payloadStr = &s
			// Try to extract ID from payload
			recordID = extractID(data)
			if userID == nil {
				userID = extractUserID(data)
			}
		}
	}

	ipAddress := c.ClientIP()
	endpoint := c.Request.URL.Path

	// Run in background — non-blocking
	go func(userID *string, method, module, endpoint string, recordID *string, payloadStr *string, ipAddress string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		query := `INSERT INTO activity_logs (user_id, method, module, endpoint, record_id, payload, ip_address) VALUES (?, ?, ?, ?, ?, ?, ?)`
		DB.WithContext(ctx).Exec(query, userID, method, module, endpoint, recordID, payloadStr, ipAddress)
	}(userID, method, module, endpoint, recordID, payloadStr, ipAddress)
}

// extractID tries to get "id" or "Id" from a struct/map via JSON marshaling.
func extractID(data any) *string {
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	if id, ok := m["id"]; ok {
		s := fmt.Sprintf("%v", id)
		if _, err := uuid.Parse(s); err == nil {
			return &s
		}
	}
	return nil
}

// extractUserID tries to get "user_id" or "created_by" from a struct/map.
func extractUserID(data any) *string {
	b, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	for _, key := range []string{"user_id", "created_by"} {
		if val, ok := m[key]; ok {
			s := fmt.Sprintf("%v", val)
			if _, err := uuid.Parse(s); err == nil {
				return &s
			}
		}
	}
	return nil
}

func handleValidationError(c *gin.Context, err error) bool {
	errMsg := err.Error()
	if strings.Contains(errMsg, "Field validation") {
		JSON(c, http.StatusBadRequest, nil, formatValidationError(errMsg))
		return true
	}
	return false
}

func formatValidationError(errStr string) string {
	errStr = strings.ReplaceAll(errStr, "Key: '", "")
	errStr = strings.ReplaceAll(errStr, "' Error:", ", ")
	errStr = strings.ReplaceAll(errStr, "Field validation for '", "")
	errStr = strings.ReplaceAll(errStr, "' failed on the '", ": ")
	errStr = strings.ReplaceAll(errStr, "' tag", "")

	re := regexp.MustCompile(`\w+Request\.`)
	errStr = re.ReplaceAllString(errStr, "")

	return errStr
}
