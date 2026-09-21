package auth_test

import (
	"backend-go/internal/auth"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	authmocks "backend-go/internal/auth/mocks"
)

const testPassword = "password123"

func mustHash(t *testing.T) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.MinCost)
	require.NoError(t, err)
	return string(hash)
}

type AuthSuite struct {
	suite.Suite
	DB      *gorm.DB
	mock    sqlmock.Sqlmock
	repo    *authmocks.MockRepositoryInterface
	service auth.ServiceInterface
}

func (s *AuthSuite) SetupTest() {
	db, sqlMock, err := sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.mock = sqlMock
	ctrl := gomock.NewController(s.T())
	s.repo = authmocks.NewMockRepositoryInterface(ctrl)
	s.service = auth.NewService(s.DB, s.repo)
}

func (s *AuthSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

func (s *AuthSuite) TestLogin_Success() {
	userID := uuid.New()
	tenantID := uuid.New()
	now := time.Now()

	user := auth.UserRow{
		ID:           userID,
		TenantID:     &tenantID,
		Name:         "Store Admin",
		Email:        "admin@laku.id",
		PasswordHash: mustHash(s.T()),
		Status:       "active",
		CreatedAt:    now,
	}

	s.repo.EXPECT().GetUserByEmail(gomock.Any(), "admin@laku.id").Return(user, nil)
	s.repo.EXPECT().UpdateLastLogin(gomock.Any(), userID, gomock.Any()).Return(nil)
	s.repo.EXPECT().GetUserRoles(gomock.Any(), userID).Return([]auth.RoleRow{
		{ID: uuid.New(), TenantID: &tenantID, Name: "admin", Permissions: 511},
	}, nil).Times(2)
	s.repo.EXPECT().CreateRefreshToken(gomock.Any(), gomock.Any()).Return(nil)

	resp, err := s.service.Login(context.Background(), &auth.LoginRequest{
		Email:    "Admin@Laku.Id",
		Password: testPassword,
	})
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp.AccessToken)
	require.NotEmpty(s.T(), resp.RefreshToken)
	require.Equal(s.T(), "Store Admin", resp.User.Name)
	require.Equal(s.T(), "admin", resp.User.Role)
	require.True(s.T(), resp.User.IsAdmin)
}

func (s *AuthSuite) TestLogin_WrongPassword() {
	user := auth.UserRow{
		ID:           uuid.New(),
		Name:         "Store Admin",
		Email:        "admin@laku.id",
		PasswordHash: mustHash(s.T()),
		Status:       "active",
	}

	s.repo.EXPECT().GetUserByEmail(gomock.Any(), "admin@laku.id").Return(user, nil)

	_, err := s.service.Login(context.Background(), &auth.LoginRequest{
		Email:    "admin@laku.id",
		Password: "wrong-password",
	})
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "email atau password salah")
}

func (s *AuthSuite) TestLogin_UserNotFound() {
	s.repo.EXPECT().GetUserByEmail(gomock.Any(), "ghost@laku.id").
		Return(auth.UserRow{}, gorm.ErrRecordNotFound)

	_, err := s.service.Login(context.Background(), &auth.LoginRequest{
		Email:    "ghost@laku.id",
		Password: testPassword,
	})
	require.Error(s.T(), err)
}

func (s *AuthSuite) TestLogin_InactiveUser() {
	user := auth.UserRow{
		ID:           uuid.New(),
		Name:         "Store Admin",
		Email:        "admin@laku.id",
		PasswordHash: mustHash(s.T()),
		Status:       "inactive",
	}

	s.repo.EXPECT().GetUserByEmail(gomock.Any(), "admin@laku.id").Return(user, nil)

	_, err := s.service.Login(context.Background(), &auth.LoginRequest{
		Email:    "admin@laku.id",
		Password: testPassword,
	})
	require.Error(s.T(), err)
}

func (s *AuthSuite) TestRefresh_RotatesToken() {
	userID := uuid.New()
	tenantID := uuid.New()
	now := time.Now()

	s.repo.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return(auth.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: "some-hash",
		ExpiresAt: now.Add(24 * time.Hour),
	}, nil)

	s.repo.EXPECT().GetUserByID(gomock.Any(), userID).Return(auth.UserRow{
		ID:        userID,
		TenantID:  &tenantID,
		Name:      "Store Admin",
		Email:     "admin@laku.id",
		Status:    "active",
		CreatedAt: now,
	}, nil)

	s.repo.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any()).Return(nil)
	s.repo.EXPECT().GetUserRoles(gomock.Any(), userID).Return([]auth.RoleRow{
		{Name: "admin", Permissions: 511},
	}, nil).Times(2)
	s.repo.EXPECT().CreateRefreshToken(gomock.Any(), gomock.Any()).Return(nil)

	resp, err := s.service.Refresh(context.Background(), &auth.RefreshRequest{RefreshToken: "valid-token"})
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp.AccessToken)
	require.NotEmpty(s.T(), resp.RefreshToken)
}

func (s *AuthSuite) TestRefresh_ExpiredToken() {
	s.repo.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return(auth.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		TokenHash: "some-hash",
		ExpiresAt: time.Now().Add(-time.Hour),
	}, nil)

	_, err := s.service.Refresh(context.Background(), &auth.RefreshRequest{RefreshToken: "expired-token"})
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "kedaluwarsa")
}

func (s *AuthSuite) TestLogout_RevokesToken() {
	s.repo.EXPECT().GetRefreshToken(gomock.Any(), gomock.Any()).Return(auth.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		TokenHash: "some-hash",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil)
	s.repo.EXPECT().RevokeRefreshToken(gomock.Any(), gomock.Any()).Return(nil)

	err := s.service.Logout(context.Background(), &auth.RefreshRequest{RefreshToken: "valid-token"})
	require.NoError(s.T(), err)
}
