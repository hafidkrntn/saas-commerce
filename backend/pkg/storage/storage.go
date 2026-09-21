package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	Delete(ctx context.Context, path string) error
	GetPresignedURL(ctx context.Context, path string, expiry time.Duration) (string, error)
	Upload(ctx context.Context, path string, reader io.Reader, size int64, contentType string) (string, error)
	UploadMultipart(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
	ValidateFileSafe(filename string) error
}

// Client wraps MinIO/S3-compatible object storage operations.
type MinioStorage struct {
	bucket string
	log    *logrus.Logger
	minio  *minio.Client
}

// Options holds storage connection parameters.
type Options struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// Init reads environment variables and initializes the global storage client.
// If S3_ENDPOINT is not set, skips without error.
func NewMinioStorage(log *logrus.Logger) (Storage, error) {
	opts := Options{
		Endpoint:  os.Getenv("MINIO_ENDPOINT"),
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		Bucket:    os.Getenv("MINIO_BUCKET"),
		UseSSL:    os.Getenv("MINIO_USE_SSL") == "true",
	}

	if opts.Endpoint == "" {
		log.Info("MINIO_ENDPOINT not set, skipping object storage")
		return nil, nil
	}

	client, err := minio.New(opts.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.AccessKey, opts.SecretKey, ""),
		Secure: opts.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("storage: failed to connect: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Ensure bucket exists
	exists, err := client.BucketExists(ctx, opts.Bucket)
	if err != nil {
		return nil, fmt.Errorf("storage: bucket check failed: %w", err)
	}

	if !exists {
		if err := client.MakeBucket(ctx, opts.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("storage: bucket create failed: %w", err)
		}
	}

	log.Info("minio storage connected")

	return &MinioStorage{
		bucket: opts.Bucket,
		log:    log,
		minio:  client,
	}, nil
}

// Delete removes an object by its full path (bucket/key).
func (c *MinioStorage) Delete(ctx context.Context, path string) error {
	if path == "" {
		return nil
	}

	// Strip bucket prefix if present
	prefix := c.bucket + "/"
	key, _ := strings.CutPrefix(path, prefix)

	return c.minio.RemoveObject(ctx, c.bucket, key, minio.RemoveObjectOptions{})
}

// GetPresignedURL generates a time-limited download URL.
func (c *MinioStorage) GetPresignedURL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	url, err := c.minio.PresignedGetObject(ctx, c.bucket, path, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("storage: presigned url failed: %w", err)
	}
	return url.String(), nil
}

// Upload stores a file and returns its object path.
func (c *MinioStorage) Upload(ctx context.Context, path string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := c.minio.PutObject(ctx, c.bucket, path, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("storage: upload failed: %w", err)
	}

	return fmt.Sprintf("%s/%s", c.bucket, path), nil
}

// UploadMultipart uploads a multipart file header with an auto-generated path.
// Rejects dangerous file extensions (executables, scripts, etc.)
func (c *MinioStorage) UploadMultipart(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	if err := c.ValidateFileSafe(file.Filename); err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("storage: failed to open file: %w", err)
	}
	defer func() { _ = src.Close() }()

	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Format("20060102-150405")
	cleanName := strings.ReplaceAll(strings.TrimSuffix(file.Filename, ext), " ", "_")
	objectName := fmt.Sprintf("%s/%s-%s%s", folder, timestamp, cleanName, ext)

	return c.Upload(ctx, objectName, src, file.Size, file.Header.Get("Content-Type"))
}

// =============================================================================
// File Safety Validation
// =============================================================================

// dangerousExtensions contains file extensions that are considered dangerous
// (executables, scripts, system files) and should never be uploaded.
var dangerousExtensions = map[string]bool{
	// Executables
	".exe": true, ".msi": true, ".dll": true, ".com": true, ".scr": true,
	".pif": true, ".app": true, ".dmg": true, ".pkg": true, ".deb": true,
	".rpm": true, ".bin": true, ".elf": true,

	// Scripts
	".bat": true, ".cmd": true, ".sh": true, ".bash": true, ".ps1": true,
	".vbs": true, ".vbe": true, ".js": true, ".jse": true, ".wsf": true,
	".wsh": true, ".py": true, ".rb": true, ".pl": true, ".php": true,
	".cgi": true, ".asp": true, ".aspx": true, ".jsx": true, ".ts": true,

	// Archives (can contain malware)
	".jar": true, ".war": true, ".ear": true,

	// System/config
	".sys": true, ".drv": true, ".inf": true, ".reg": true, ".lnk": true,
	".iso": true, ".img": true,

	// Office macros
	".docm": true, ".xlsm": true, ".pptm": true, ".dotm": true,

	// Others
	".hta": true, ".cpl": true, ".msp": true, ".mst": true,
}

// ValidateFileSafe checks if a file has a dangerous extension.
// Returns an error if the file extension is blacklisted.
func (c *MinioStorage) ValidateFileSafe(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return fmt.Errorf("storage: file tanpa ekstensi tidak diizinkan")
	}
	if dangerousExtensions[ext] {
		return fmt.Errorf("storage: file dengan ekstensi '%s' tidak diizinkan (file berbahaya)", ext)
	}
	return nil
}
