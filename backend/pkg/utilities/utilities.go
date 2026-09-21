package utilities

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =============================================================================
// UUID Helpers
// =============================================================================

// ParseUUID parses a string to uuid.UUID with a friendly error message.
func ParseUUID(id string) (uuid.UUID, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("id tidak valid: %s", id)
	}
	return uid, nil
}

// =============================================================================
// String Helpers
// =============================================================================

// IsEmpty returns true if the string pointer is nil or blank.
func IsEmpty(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}

// Contains checks if a string exists in a slice.
func Contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

// =============================================================================
// Number Formatting
// =============================================================================

// FormatNumber formats an integer with dot separators (Indonesian style).
// Example: 1000000 → "1.000.000"
func FormatNumber(n float64) string {
	intPart := int64(n)
	num := strconv.FormatInt(intPart, 10)

	var sb strings.Builder
	sb.Grow(len(num) + len(num)/3)

	for i, c := range num {
		pos := len(num) - i
		if i > 0 && pos%3 == 0 {
			sb.WriteByte('.')
		}
		sb.WriteRune(c)
	}

	return sb.String()
}

// =============================================================================
// Date/Time Helpers
// =============================================================================

// ParseDate tries to parse a date string in multiple common formats.
func ParseDate(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", s)
}

// SameDate returns true if two times fall on the same calendar date.
func SameDate(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

// =============================================================================
// File Validation
// =============================================================================

var (
	allowedImageExts = map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	}
	allowedPDFExts = map[string]bool{
		".pdf": true,
	}
)

// ValidateFileSize checks if a file exceeds the given size limit in MB.
func ValidateFileSize(fileSize int64, maxSizeMB int) error {
	maxBytes := int64(maxSizeMB) * 1024 * 1024
	if fileSize > maxBytes {
		return fmt.Errorf("ukuran file melebihi batas %d MB", maxSizeMB)
	}
	return nil
}

// ValidateImageFile checks if the file has an allowed image extension.
func ValidateImageFile(filename string) error {
	if !allowedImageExts[strings.ToLower(filepath.Ext(filename))] {
		return fmt.Errorf("file harus berupa gambar (jpg, jpeg, png, webp): %s", filename)
	}
	return nil
}

// ValidatePDFFile checks if the file has a .pdf extension.
func ValidatePDFFile(filename string) error {
	if !allowedPDFExts[strings.ToLower(filepath.Ext(filename))] {
		return fmt.Errorf("file harus berupa PDF: %s", filename)
	}
	return nil
}

// ValidateUploadFile validates a multipart file for allowed types and size (max 5MB).
func ValidateUploadFile(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExts[ext] && !allowedPDFExts[ext] {
		return fmt.Errorf("file harus berupa gambar atau PDF: %s", file.Filename)
	}
	return ValidateFileSize(file.Size, 5)
}

// =============================================================================
// Transaction Number Generator
// =============================================================================

// GenerateTransactionNumber creates a sequential transaction number with format:
// PREFIX/YYMM/00001 (in-memory, you provide the last number)
func GenerateTransactionNumber(prefix string, lastNumber string) string {
	datePart := time.Now().Format("0601")
	fullPrefix := prefix + "/" + datePart + "/"

	nextNum := 1
	if lastNumber != "" && strings.HasPrefix(lastNumber, fullPrefix) {
		numStr := lastNumber[len(lastNumber)-5:]
		if n, err := strconv.Atoi(numStr); err == nil {
			nextNum = n + 1
		}
	}

	return fmt.Sprintf("%s%05d", fullPrefix, nextNum)
}

// GenerateNumberTransactions queries the DB for the last number and generates the next one.
// Format: PREFIX/YYMM/00001
//
// Params:
//   - db: gorm.DB instance (use tx if inside a transaction)
//   - table: table name (e.g. "orders", "invoices")
//   - column: column name that stores the number (e.g. "order_number", "invoice_number")
//   - prefix: prefix string (e.g. "ORD", "INV", "SQ")
//
// Example:
//
//	number := utilities.GenerateNumberTransactions(tx, "orders", "order_number", "ORD")
//	// Result: "ORD/2506/00001"
func GenerateNumberTransactions(db *gorm.DB, table string, column string, prefix string) string {
	datePart := time.Now().Format("0601")
	fullPrefix := prefix + "/" + datePart + "/"
	likePattern := fullPrefix + "%"

	var lastNumber string
	db.Table(table).
		Select(column).
		Where(column+" LIKE ?", likePattern).
		Order(column+" DESC").
		Limit(1).
		Pluck(column, &lastNumber)

	nextNum := 1
	if lastNumber != "" {
		numStr := lastNumber[len(lastNumber)-5:]
		if n, err := strconv.Atoi(numStr); err == nil {
			nextNum = n + 1
		}
	}

	return fmt.Sprintf("%s%05d", fullPrefix, nextNum)
}

// GenerateNumberTransactionsScoped is the tenant-scoped variant of
// GenerateNumberTransactions. The sequence restarts per tenant.
func GenerateNumberTransactionsScoped(db *gorm.DB, tenantID uuid.UUID, table string, column string, prefix string) string {
	datePart := time.Now().Format("0601")
	fullPrefix := prefix + "/" + datePart + "/"
	likePattern := fullPrefix + "%"

	var lastNumber string
	db.Table(table).
		Select(column).
		Where(column+" LIKE ? AND tenant_id = ?", likePattern, tenantID).
		Order(column+" DESC").
		Limit(1).
		Pluck(column, &lastNumber)

	nextNum := 1
	if lastNumber != "" {
		numStr := lastNumber[len(lastNumber)-5:]
		if n, err := strconv.Atoi(numStr); err == nil {
			nextNum = n + 1
		}
	}

	return fmt.Sprintf("%s%05d", fullPrefix, nextNum)
}
