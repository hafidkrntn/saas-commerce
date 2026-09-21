package token

import (
	"backend-go/internal/entities"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// jwtKey and apiKey are lazily loaded so that .env (loaded in main) is
// respected. Package-level init would run before godotenv.Load().
var (
	jwtKeyOnce sync.Once
	jwtKey     []byte

	apiKeyOnce sync.Once
	apiKey     string
)

func getSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "IuXPToVgeRnMgaFXQQ9m2ApHSCtLrOrA" // jangan dipakai di production
	}
	return secret
}

func getAPIKey() string {
	key := os.Getenv("API_KEY")
	if key == "" {
		key = "RcFKQ_xvKBCe6eL9ekf4Z3HDYx0jCgkx" // jangan dipakai di production
	}
	return key
}

func jwtSecret() []byte {
	jwtKeyOnce.Do(func() {
		jwtKey = []byte(getSecret())
	})
	return jwtKey
}

func apiKeyValue() string {
	apiKeyOnce.Do(func() {
		apiKey = getAPIKey()
	})
	return apiKey
}

// --- Password Hashing ---
func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}

// --- Random Token Generation ---

// GenerateRandomRefreshToken generates a cryptographically secure random string
func GenerateRandomRefreshToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// --- JWT Token Generation ---

// GenerateToken - deprecated, use GenerateAccessToken instead
func GenerateToken(session entities.SessionToken) (string, error) {
	return GenerateAccessToken(session, 24*time.Hour)
}

// GenerateAccessToken generates a short-lived access token (JWT)
func GenerateAccessToken(session entities.SessionToken, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{}

	// Use reflection to automatically get all struct fields
	val := reflect.ValueOf(session)
	typ := reflect.TypeOf(session)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		claims[field.Name] = value
	}

	// Add token metadata
	claims["exp"] = time.Now().Add(ttl).Unix()
	claims["iat"] = time.Now().Unix()
	claims["type"] = "access_token"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

// --- JWT Token Verification ---

// VerifyToken - for access tokens
func VerifyToken(c *gin.Context) (jwt.MapClaims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header missing")
	}

	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		return nil, errors.New("invalid authorization format")
	}
	tokenString := splitToken[1]

	return verifyAccessTokenString(tokenString)
}

// verifyAccessTokenString verifies JWT access token
func verifyAccessTokenString(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token type
		if tokenType, exists := claims["type"].(string); !exists || tokenType != "access_token" {
			return nil, errors.New("invalid token type")
		}

		// Check expiration
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, errors.New("invalid expiration in token")
		}
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			return nil, errors.New("token has expired")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// --- API Key Verification ---
func VerifyAPIKey(c *gin.Context) error {
	key := c.GetHeader("X-API-KEY")
	if key == "" {
		return errors.New("api key missing")
	}
	if key != apiKeyValue() {
		return errors.New("invalid api key")
	}
	return nil
}

// --- Helper Functions ---

// ExtractUserIDFromToken extracts user ID from token claims
func ExtractUserIDFromToken(claims jwt.MapClaims) (string, error) {
	userID, ok := claims["UserId"].(string)
	if !ok {
		// Try alternative field names
		if uid, exists := claims["user_id"].(string); exists {
			return uid, nil
		}
		return "", errors.New("user_id not found in token")
	}
	return userID, nil
}

// IsTokenExpired checks if token is expired without full verification
func IsTokenExpired(tokenString string) bool {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return true
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(exp), 0)
			return time.Now().After(expTime)
		}
	}
	return true
}

// GetTokenClaims gets claims from token without verification (use carefully)
func GetTokenClaims(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid claims")
}
