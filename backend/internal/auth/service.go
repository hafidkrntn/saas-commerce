package auth

import (
	"backend-go/internal/entities"
	"backend-go/pkg/apperror"
	"backend-go/pkg/token"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Token lifetimes
const (
	accessTokenTTL  = 24 * time.Hour
	refreshTokenTTL = 30 * 24 * time.Hour
)

type ServiceInterface interface {
	Login(ctx context.Context, req *LoginRequest) (LoginResponse, error)
	Refresh(ctx context.Context, req *RefreshRequest) (LoginResponse, error)
	Logout(ctx context.Context, req *RefreshRequest) error
	Me(ctx context.Context, userID string) (LoginResponse, error)
}

type Service struct {
	db   *gorm.DB
	repo RepositoryInterface
}

func NewService(db *gorm.DB, repo RepositoryInterface) ServiceInterface {
	return &Service{db: db, repo: repo}
}

func (s *Service) Login(ctx context.Context, req *LoginRequest) (LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.repo.GetUserByEmail(s.db.WithContext(ctx), email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResponse{}, apperror.Unauthorized("email atau password salah", nil)
		}
		return LoginResponse{}, apperror.Internal("gagal memeriksa user", err)
	}

	if user.Status != "active" {
		return LoginResponse{}, apperror.Forbidden("akun tidak aktif", nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return LoginResponse{}, apperror.Unauthorized("email atau password salah", nil)
	}

	now := time.Now()
	if err := s.repo.UpdateLastLogin(s.db.WithContext(ctx), user.ID, now); err != nil {
		return LoginResponse{}, apperror.Internal("gagal memperbarui last login", err)
	}

	user.LastLoginAt = &now
	return s.issueTokens(ctx, user)
}

func (s *Service) Refresh(ctx context.Context, req *RefreshRequest) (LoginResponse, error) {
	tokenHash := token.HashToken(req.RefreshToken)

	refreshToken, err := s.repo.GetRefreshToken(s.db.WithContext(ctx), tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResponse{}, apperror.Unauthorized("refresh token tidak valid", nil)
		}
		return LoginResponse{}, apperror.Internal("gagal memeriksa refresh token", err)
	}

	if refreshToken.RevokedAt != nil {
		return LoginResponse{}, apperror.Unauthorized("refresh token sudah tidak berlaku", nil)
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return LoginResponse{}, apperror.Unauthorized("refresh token kedaluwarsa", nil)
	}

	user, err := s.repo.GetUserByID(s.db.WithContext(ctx), refreshToken.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResponse{}, apperror.Unauthorized("user tidak ditemukan", nil)
		}
		return LoginResponse{}, apperror.Internal("gagal mengambil user", err)
	}

	if user.Status != "active" {
		return LoginResponse{}, apperror.Forbidden("akun tidak aktif", nil)
	}

	// Rotate: revoke old, issue new pair
	if err := s.repo.RevokeRefreshToken(s.db.WithContext(ctx), refreshToken.ID); err != nil {
		return LoginResponse{}, apperror.Internal("gagal mencabut refresh token lama", err)
	}

	return s.issueTokens(ctx, user)
}

func (s *Service) Logout(ctx context.Context, req *RefreshRequest) error {
	tokenHash := token.HashToken(req.RefreshToken)

	refreshToken, err := s.repo.GetRefreshToken(s.db.WithContext(ctx), tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return apperror.Internal("gagal memeriksa refresh token", err)
	}

	return s.repo.RevokeRefreshToken(s.db.WithContext(ctx), refreshToken.ID)
}

// Me returns a fresh token pair for the authenticated user (session refresh).
func (s *Service) Me(ctx context.Context, userID string) (LoginResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return LoginResponse{}, apperror.BadRequest("user id tidak valid", err)
	}

	user, err := s.repo.GetUserByID(s.db.WithContext(ctx), uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResponse{}, apperror.NotFound("user tidak ditemukan", err)
		}
		return LoginResponse{}, apperror.Internal("gagal mengambil user", err)
	}

	return s.buildUserResponse(ctx, user)
}

// =============================================================================
// Token issuance
// =============================================================================

// issueTokens generates a fresh access + refresh token pair for the user.
func (s *Service) issueTokens(ctx context.Context, user UserRow) (LoginResponse, error) {
	roles, err := s.repo.GetUserRoles(s.db.WithContext(ctx), user.ID)
	if err != nil {
		return LoginResponse{}, apperror.Internal("gagal mengambil role user", err)
	}

	role, isAdmin := resolveRole(roles)

	tenantID := ""
	userGroupID := ""
	if user.TenantID != nil {
		tenantID = user.TenantID.String()
		userGroupID = tenantID
	}

	session := entities.SessionToken{
		SessionId:       uuid.NewString(),
		UserId:          user.ID.String(),
		UserGroupId:     userGroupID,
		TenantId:        tenantID,
		Role:            role,
		IsAdministrator: isAdmin,
		Application:     "saas-ecommerce",
	}

	accessToken, err := token.GenerateAccessToken(session, accessTokenTTL)
	if err != nil {
		return LoginResponse{}, apperror.Internal("gagal membuat access token", err)
	}

	refreshToken, err := token.GenerateRandomRefreshToken()
	if err != nil {
		return LoginResponse{}, apperror.Internal("gagal membuat refresh token", err)
	}

	if err := s.repo.CreateRefreshToken(s.db.WithContext(ctx), &RefreshToken{
		UserID:    user.ID,
		TokenHash: token.HashToken(refreshToken),
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}); err != nil {
		return LoginResponse{}, apperror.Internal("gagal menyimpan refresh token", err)
	}

	resp, err := s.buildUserResponse(ctx, user)
	if err != nil {
		return LoginResponse{}, err
	}

	resp.AccessToken = accessToken
	resp.TokenType = "Bearer"
	resp.ExpiresIn = int64(accessTokenTTL.Seconds())
	resp.RefreshToken = refreshToken
	return resp, nil
}

func (s *Service) buildUserResponse(ctx context.Context, user UserRow) (LoginResponse, error) {
	roles, err := s.repo.GetUserRoles(s.db.WithContext(ctx), user.ID)
	if err != nil {
		return LoginResponse{}, apperror.Internal("gagal mengambil role user", err)
	}

	role, isAdmin := resolveRole(roles)

	return LoginResponse{
		User: UserResponse{
			ID:        user.ID,
			TenantID:  user.TenantID,
			Name:      user.Name,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Role:      role,
			IsAdmin:   isAdmin,
			LastLogin: user.LastLoginAt,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// resolveRole picks the user's primary role by permission weight.
func resolveRole(roles []RoleRow) (string, bool) {
	if len(roles) == 0 {
		return "", false
	}

	best := roles[0]
	for _, r := range roles[1:] {
		if r.Permissions > best.Permissions {
			best = r
		}
	}

	isAdmin := best.Name == "admin" || best.Name == "platform-admin"
	return best.Name, isAdmin
}
