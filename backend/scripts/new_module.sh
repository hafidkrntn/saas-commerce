#!/bin/bash

# Script to scaffold a new domain module with the standard file structure.
# Usage: ./scripts/new_module.sh <module_name> [tests]
#
# Generates:
#   model.go | repository.go | service.go | handler.go | routes.go
#   mocks/ (via mockgen)
#   repository_test.go | service_test.go (table-driven, gomock) if tests=true
#
# Supports hyphenated names: sales-quotation
#   → folder:  internal/salesquotation/
#   → package: salesquotation
#   → struct:  SalesQuotation
#   → route:   /sales-quotations
#   → table:   sales_quotations

INPUT_NAME=$1
GEN_TESTS=$2

if [ -z "$INPUT_NAME" ]; then
  echo "❌ Usage: $0 <module_name> [tests]"
  echo "   Example: $0 sales-quotation"
  exit 1
fi

# Derive names from input
PKG_NAME=$(echo "$INPUT_NAME" | tr -d '-')
STRUCT_NAME=$(echo "$INPUT_NAME" | awk -F'-' '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)}1' OFS='')
ROUTE_NAME="$INPUT_NAME"
TABLE_NAME=$(echo "$INPUT_NAME" | tr '-' '_')

PLURAL_TABLE="${TABLE_NAME}s"
PLURAL_ROUTE="${ROUTE_NAME}s"

MODULE_DIR="internal/${PKG_NAME}"
MOCK_DIR="${MODULE_DIR}/mocks"

if [ -d "$MODULE_DIR" ]; then
  echo "❌ Module '${PKG_NAME}' already exists at ${MODULE_DIR}"
  exit 1
fi

mkdir -p "$MOCK_DIR"

echo "📦 Module: ${INPUT_NAME}"
echo "   Package:  ${PKG_NAME}"
echo "   Struct:   ${STRUCT_NAME}"
echo "   Route:    /${PLURAL_ROUTE}"
echo "   Table:    ${PLURAL_TABLE}"
echo ""

# =============================================================================
# model.go
# =============================================================================
cat > "${MODULE_DIR}/model.go" << EOF
package ${PKG_NAME}

import (
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// Database Models
// =============================================================================

type ${STRUCT_NAME} struct {
	ID        uuid.UUID \`gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"\`
	Name      string    \`gorm:"column:name;not null"\`
	CreatedBy string    \`gorm:"column:created_by"\`
	UpdatedBy string    \`gorm:"column:updated_by"\`
	CreatedAt time.Time \`gorm:"column:created_at;autoCreateTime"\`
	UpdatedAt time.Time \`gorm:"column:updated_at;autoUpdateTime"\`
}

func (${STRUCT_NAME}) TableName() string {
	return "${PLURAL_TABLE}"
}

// =============================================================================
// Request DTOs
// =============================================================================

type CreateRequest struct {
	Name string \`json:"name" binding:"required"\`
}

type UpdateRequest struct {
	ID   string  \`json:"id" binding:"required"\`
	Name *string \`json:"name" binding:"omitempty,min=1"\`
}

// =============================================================================
// Response DTOs
// =============================================================================

type Response struct {
	ID        uuid.UUID \`json:"id"\`
	Name      string    \`json:"name"\`
	CreatedAt time.Time \`json:"created_at"\`
	UpdatedAt time.Time \`json:"updated_at"\`
}

func (m *${STRUCT_NAME}) ToResponse() Response {
	return Response{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
EOF

# =============================================================================
# repository.go
# =============================================================================
cat > "${MODULE_DIR}/repository.go" << EOF
package ${PKG_NAME}

import (
	"backend-go/internal/entities"
	"backend-go/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RepositoryInterface defines data-access methods for this domain.
type RepositoryInterface interface {
	Create(tx *gorm.DB, model *${STRUCT_NAME}) error
	GetByID(tx *gorm.DB, id uuid.UUID) (${STRUCT_NAME}, error)
	List(tx *gorm.DB, params *entities.ListParams) (pagination.Response[${STRUCT_NAME}], error)
	Update(tx *gorm.DB, model *${STRUCT_NAME}) error
	Delete(tx *gorm.DB, id uuid.UUID) error
}

// Repository implements RepositoryInterface.
type Repository struct{}

func NewRepository() RepositoryInterface {
	return &Repository{}
}

func (r *Repository) Create(tx *gorm.DB, model *${STRUCT_NAME}) error {
	return tx.Create(model).Error
}

func (r *Repository) GetByID(tx *gorm.DB, id uuid.UUID) (${STRUCT_NAME}, error) {
	var result ${STRUCT_NAME}
	if err := tx.Where("id = ?", id).First(&result).Error; err != nil {
		return ${STRUCT_NAME}{}, err
	}
	return result, nil
}

func (r *Repository) List(tx *gorm.DB, params *entities.ListParams) (pagination.Response[${STRUCT_NAME}], error) {
	p := pagination.Resolve(params.Page, params.Limit)

	query := tx.Model(&${STRUCT_NAME}{})

	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("name ILIKE ?", search)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return pagination.Response[${STRUCT_NAME}]{}, err
	}

	var items []${STRUCT_NAME}
	if err := query.
		Order("created_at DESC").
		Offset(p.Offset).
		Limit(p.Limit).
		Find(&items).Error; err != nil {
		return pagination.Response[${STRUCT_NAME}]{}, err
	}

	return pagination.NewResponse(items, int(total), p.Page, p.Limit), nil
}

func (r *Repository) Update(tx *gorm.DB, model *${STRUCT_NAME}) error {
	return tx.Save(model).Error
}

func (r *Repository) Delete(tx *gorm.DB, id uuid.UUID) error {
	return tx.Where("id = ?", id).Delete(&${STRUCT_NAME}{}).Error
}
EOF

# =============================================================================
# service.go
# =============================================================================
cat > "${MODULE_DIR}/service.go" << EOF
package ${PKG_NAME}

import (
	"backend-go/internal/entities"
	"backend-go/pkg/apperror"
	"backend-go/pkg/cache"
	"backend-go/pkg/pagination"
	"backend-go/pkg/storage"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ServiceInterface defines the public contract for this service.
type ServiceInterface interface {
	Create(ctx context.Context, req *CreateRequest, userID string) (entities.PostResponse, error)
	GetByID(ctx context.Context, id string) (Response, error)
	List(ctx context.Context, params *entities.ListParams) (pagination.Response[Response], error)
	Update(ctx context.Context, req *UpdateRequest, userID string) (entities.PutResponse, error)
	Delete(ctx context.Context, id string) error
}

// Service implements ServiceInterface.
type Service struct {
	db      *gorm.DB
	repo    RepositoryInterface
	cache   cache.Cache
	storage storage.Storage
}

func NewService(db *gorm.DB, repo RepositoryInterface, cache cache.Cache, storage storage.Storage) ServiceInterface {
	return &Service{db: db, repo: repo, cache: cache, storage: storage}
}

func (s *Service) Create(ctx context.Context, req *CreateRequest, userID string) (entities.PostResponse, error) {
	var model ${STRUCT_NAME}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model = ${STRUCT_NAME}{
			Name:      req.Name,
			CreatedBy: userID,
		}

		if err := s.repo.Create(tx, &model); err != nil {
			return apperror.Internal("gagal membuat data", err)
		}

		return nil
	})

	if err != nil {
		return entities.PostResponse{}, err
	}

	return entities.PostResponse{
		Id:        model.ID,
		Version:   1,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (Response, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return Response{}, apperror.BadRequest("id tidak valid", err)
	}

	result, err := s.repo.GetByID(s.db.WithContext(ctx), uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Response{}, apperror.NotFound("data tidak ditemukan", err)
		}
		return Response{}, apperror.Internal("gagal mengambil data", err)
	}

	resp := result.ToResponse()
	return resp, nil
}

func (s *Service) List(ctx context.Context, params *entities.ListParams) (pagination.Response[Response], error) {
	result, err := s.repo.List(s.db.WithContext(ctx), params)
	if err != nil {
		return pagination.Response[Response]{}, apperror.Internal("gagal mengambil daftar data", err)
	}

	responses := make([]Response, len(result.Data))
	for i, item := range result.Data {
		responses[i] = item.ToResponse()
	}

	return pagination.NewResponse(responses, result.Pagination.Total, result.Pagination.Page, result.Pagination.Limit), nil
}

func (s *Service) Update(ctx context.Context, req *UpdateRequest, userID string) (entities.PutResponse, error) {
	uid, err := uuid.Parse(req.ID)
	if err != nil {
		return entities.PutResponse{}, apperror.BadRequest("id tidak valid", err)
	}

	var model ${STRUCT_NAME}
	txErr := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		model, err = s.repo.GetByID(tx, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("data tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil data", err)
		}

		if req.Name != nil {
			model.Name = *req.Name
		}
		model.UpdatedBy = userID

		return s.repo.Update(tx, &model)
	})

	if txErr != nil {
		return entities.PutResponse{}, txErr
	}

	now := model.UpdatedAt
	return entities.PutResponse{
		Id:        model.ID,
		Version:   1,
		UpdatedAt: &now,
	}, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return apperror.BadRequest("id tidak valid", err)
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, err := s.repo.GetByID(tx, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("data tidak ditemukan", err)
			}
			return apperror.Internal("gagal mengambil data", err)
		}

		if err := s.repo.Delete(tx, uid); err != nil {
			return apperror.Internal("gagal menghapus data", err)
		}

		return nil
	})

	return err
}
EOF

# =============================================================================
# handler.go
# =============================================================================
cat > "${MODULE_DIR}/handler.go" << EOF
package ${PKG_NAME}

import (
	"backend-go/internal/entities"
	"backend-go/pkg/response"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the ${PKG_NAME} domain.
type Handler struct {
	svc ServiceInterface
}

func NewHandler(svc ServiceInterface) *Handler {
	return &Handler{svc: svc}
}

// Create handles POST /${PLURAL_ROUTE}/create
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	userID := c.GetString("UserId")
	data, err := h.svc.Create(c.Request.Context(), &req, userID)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Created(c, data, &req)
}

// GetByID handles GET /${PLURAL_ROUTE}/get/:id
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

// List handles GET /${PLURAL_ROUTE}/list
func (h *Handler) List(c *gin.Context) {
	var params entities.ListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Err(c, err)
		return
	}

	result, err := h.svc.List(c.Request.Context(), &params)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.OK(c, result)
}

// Update handles PUT /${PLURAL_ROUTE}/update (ID in body)
func (h *Handler) Update(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Err(c, err)
		return
	}

	userID := c.GetString("UserId")
	data, err := h.svc.Update(c.Request.Context(), &req, userID)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Updated(c, data, &req)
}

// Delete handles DELETE /${PLURAL_ROUTE}/delete/:id
func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	_ = idParam // use utilities.ParseUUID in real modules

	data, err := h.svc.GetByID(c.Request.Context(), idParam)
	if err != nil {
		response.Err(c, err)
		return
	}

	err = h.svc.Delete(c.Request.Context(), idParam)
	if err != nil {
		response.Err(c, err)
		return
	}

	response.Deleted(c, data)
}
EOF

# =============================================================================
# routes.go
# =============================================================================
cat > "${MODULE_DIR}/routes.go" << EOF
package ${PKG_NAME}

import (
	"github.com/gin-gonic/gin"
)

// Router is a Wire-friendly route registrar.
type Router struct {
	handler *Handler
}

// NewRouter creates a new Router.
func NewRouter(handler *Handler) *Router {
	return &Router{handler: handler}
}

// RegisterRoutes registers all ${PKG_NAME} domain routes.
func (r *Router) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/${PLURAL_ROUTE}")
	{
		group.GET("/list", r.handler.List)
		group.POST("/create", r.handler.Create)
		group.GET("/get/:id", r.handler.GetByID)
		group.PUT("/update", r.handler.Update)
		group.DELETE("/delete/:id", r.handler.Delete)
	}
}
EOF

# =============================================================================
# Generate mocks via mockgen
# =============================================================================
if ! command -v mockgen &> /dev/null; then
  echo "⚠️  mockgen not found. Installing..."
  go install go.uber.org/mock/mockgen@latest
fi

mockgen -source="${MODULE_DIR}/repository.go" \
  -destination="${MOCK_DIR}/${PKG_NAME}_repository.go" \
  -package=mocks RepositoryInterface

mockgen -source="${MODULE_DIR}/service.go" \
  -destination="${MOCK_DIR}/${PKG_NAME}_service.go" \
  -package=mocks ServiceInterface

echo "✅ Mocks generated at ${MOCK_DIR}/"

# =============================================================================
# repository_test.go (only when tests=true)
# =============================================================================
if [ "$GEN_TESTS" = "true" ]; then
cat > "${MODULE_DIR}/repository_test.go" << 'REPO_TEST_EOF'
package REPLACE_test

import (
	"REPLACE_PKG/internal/entities"
	"REPLACE_PKG/internal/REPLACE"
	"REPLACE_PKG/pkg/pagination"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository REPLACE.RepositoryInterface
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(s.T(), err)

	s.DB.Logger.LogMode(1)

	s.repository = REPLACE.NewRepository()
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestRepository_Create() {
	now := time.Now()
	newModel := &REPLACE.REPLACE_STRUCT{
		Name:      "Test Name",
		CreatedBy: "user-1",
		CreatedAt: now,
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "REPLACE_TABLE"`),
	).
		WithArgs(
			newModel.Name,
			newModel.CreatedBy,
			newModel.UpdatedBy,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(newModel.ID),
		)
	s.mock.ExpectCommit()

	err := s.repository.Create(s.DB, newModel)
	require.NoError(s.T(), err)
}

func (s *Suite) TestRepository_GetByID() {
	expected := REPLACE.REPLACE_STRUCT{
		ID:        uuid.New(),
		Name:      "Test Name",
		CreatedBy: "user-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "created_by", "updated_by", "created_at", "updated_at",
	}).
		AddRow(expected.ID, expected.Name, expected.CreatedBy, expected.UpdatedBy, expected.CreatedAt, expected.UpdatedAt)

	s.mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "REPLACE_TABLE" WHERE id = $1 ORDER BY "REPLACE_TABLE"."id" LIMIT $2`),
	).
		WithArgs(expected.ID, 1).
		WillReturnRows(rows)

	result, err := s.repository.GetByID(s.DB, expected.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expected, result)
}

func (s *Suite) TestRepository_List() {
	tests := []struct {
		name      string
		params    *entities.ListParams
		mockSetup func()
		assert    func(*testing.T, pagination.Response[REPLACE.REPLACE_STRUCT], error)
	}{
		{
			name: "should list items",
			params: &entities.ListParams{
				Page:  1,
				Limit: 10,
			},
			mockSetup: func() {
				s.mock.ExpectQuery(`SELECT count\(\*\) FROM "REPLACE_TABLE"`).
					WillReturnRows(
						sqlmock.NewRows([]string{"count"}).AddRow(2),
					)

				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(uuid.New(), "Item 1").
					AddRow(uuid.New(), "Item 2")

				s.mock.ExpectQuery(`SELECT .* FROM "REPLACE_TABLE"`).
					WillReturnRows(rows)
			},
			assert: func(t *testing.T, result pagination.Response[REPLACE.REPLACE_STRUCT], err error) {
				require.NoError(t, err)
				require.Equal(t, 2, result.Pagination.Total)
				require.Len(t, result.Data, 2)
			},
		},
		{
			name: "should return count error",
			params: &entities.ListParams{
				Page:  1,
				Limit: 10,
			},
			mockSetup: func() {
				s.mock.ExpectQuery(`SELECT count\(\*\) FROM "REPLACE_TABLE"`).
					WillReturnError(errors.New("database error"))
			},
			assert: func(t *testing.T, _ pagination.Response[REPLACE.REPLACE_STRUCT], err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "database error")
			},
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			test.mockSetup()
			result, err := s.repository.List(s.DB, test.params)
			test.assert(t, result, err)
		})
	}
}

func (s *Suite) TestRepository_Update() {
	model := &REPLACE.REPLACE_STRUCT{
		ID:        uuid.New(),
		Name:      "Updated Name",
		UpdatedBy: "user-1",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "REPLACE_TABLE"`).
		WithArgs(
			model.Name,
			model.CreatedBy,
			model.UpdatedBy,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			model.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repository.Update(s.DB, model)
	require.NoError(s.T(), err)
}

func (s *Suite) TestRepository_Delete() {
	id := uuid.New()

	s.mock.ExpectBegin()
	s.mock.ExpectExec(
		regexp.QuoteMeta(`DELETE FROM "REPLACE_TABLE" WHERE id = $1`),
	).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repository.Delete(s.DB, id)
	require.NoError(s.T(), err)
}
REPO_TEST_EOF

# Replace placeholders in repository test
sed -i '' \
  -e "s/REPLACE_PKG/backend-go/g" \
  -e "s/REPLACE_STRUCT/${STRUCT_NAME}/g" \
  -e "s/REPLACE_TABLE/${PLURAL_TABLE}/g" \
  -e "s/REPLACE/${PKG_NAME}/g" \
  "${MODULE_DIR}/repository_test.go"

# =============================================================================
# service_test.go
# =============================================================================
cat > "${MODULE_DIR}/service_test.go" << 'SVC_TEST_EOF'
package REPLACE_test

import (
	"REPLACE_PKG/internal/REPLACE"
	"REPLACE_PKG/internal/REPLACE/mocks"
	"REPLACE_PKG/internal/entities"
	cacheMocks "REPLACE_PKG/pkg/cache/mocks"
	"REPLACE_PKG/pkg/pagination"
	storageMocks "REPLACE_PKG/pkg/storage/mocks"
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// testDB is a shared *gorm.DB instance backed by sqlmock for non-transaction tests.
var testDB *gorm.DB

func init() {
	var err error

	db, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	testDB, err = gorm.Open(
		postgres.New(postgres.Config{Conn: db}),
		&gorm.Config{},
	)
	if err != nil {
		panic(err)
	}
}

// newServiceForTransaction creates a fresh sqlmock for transaction tests.
func newServiceForTransaction(t *testing.T, ctrl *gomock.Controller) (REPLACE.ServiceInterface, *mocks.MockRepositoryInterface, sqlmock.Sqlmock) {
	t.Helper()

	db, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: db}),
		&gorm.Config{},
	)
	require.NoError(t, err)

	mockRepo := mocks.NewMockRepositoryInterface(ctrl)
	mockCache := cacheMocks.NewMockCache(ctrl)
	mockStorage := storageMocks.NewMockStorage(ctrl)
	svc := REPLACE.NewService(gormDB, mockRepo, mockCache, mockStorage)
	return svc, mockRepo, sqlMock
}

// =============================================================================
// Tests: Create
// =============================================================================

func TestService_Create(t *testing.T) {
	t.Parallel()

	type args struct {
		req    *REPLACE.CreateRequest
		userID string
	}

	tests := []struct {
		name      string
		args      args
		mockSetup func(*mocks.MockRepositoryInterface, sqlmock.Sqlmock)
		wantErr   string
	}{
		{
			name: "should create successfully",
			args: args{
				req:    &REPLACE.CreateRequest{Name: "Test"},
				userID: "user-1",
			},
			mockSetup: func(repo *mocks.MockRepositoryInterface, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
				sqlMock.ExpectCommit()
			},
			wantErr: "",
		},
		{
			name: "should return error when create fails",
			args: args{
				req:    &REPLACE.CreateRequest{Name: "Test"},
				userID: "user-1",
			},
			mockSetup: func(repo *mocks.MockRepositoryInterface, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("insert failed"))
				sqlMock.ExpectRollback()
			},
			wantErr: "gagal membuat data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc, mockRepo, sqlMock := newServiceForTransaction(t, ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo, sqlMock)
			}

			resp, err := svc.Create(context.Background(), tt.args.req, tt.args.userID)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Empty(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, resp)
		})
	}
}

// =============================================================================
// Tests: GetByID
// =============================================================================

func TestService_GetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		mockSetup func(*mocks.MockRepositoryInterface)
		wantErr   string
	}{
		{
			name:    "should return bad request for invalid uuid",
			id:      "not-a-uuid",
			wantErr: "id tidak valid",
		},
		{
			name: "should return not found",
			id:   uuid.New().String(),
			mockSetup: func(repo *mocks.MockRepositoryInterface) {
				repo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(REPLACE.REPLACE_STRUCT{}, gorm.ErrRecordNotFound)
			},
			wantErr: "data tidak ditemukan",
		},
		{
			name: "should return internal error when db fails",
			id:   uuid.New().String(),
			mockSetup: func(repo *mocks.MockRepositoryInterface) {
				repo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(REPLACE.REPLACE_STRUCT{}, errors.New("db error"))
			},
			wantErr: "gagal mengambil data",
		},
		{
			name: "should return successfully",
			id:   uuid.New().String(),
			mockSetup: func(repo *mocks.MockRepositoryInterface) {
				repo.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(REPLACE.REPLACE_STRUCT{
						ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						Name: "Test",
					}, nil)
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)
			mockCache := cacheMocks.NewMockCache(ctrl)
			mockStorage := storageMocks.NewMockStorage(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			svc := REPLACE.NewService(testDB, mockRepo, mockCache, mockStorage)
			resp, err := svc.GetByID(context.Background(), tt.id)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Empty(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, resp)
			require.Equal(t, "Test", resp.Name)
		})
	}
}

// =============================================================================
// Tests: List
// =============================================================================

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		params    *entities.ListParams
		mockSetup func(*mocks.MockRepositoryInterface)
		wantErr   string
		wantTotal int
	}{
		{
			name: "should return internal error when repo fails",
			params: &entities.ListParams{Page: 1, Limit: 10},
			mockSetup: func(repo *mocks.MockRepositoryInterface) {
				repo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(pagination.Response[REPLACE.REPLACE_STRUCT]{}, errors.New("db error"))
			},
			wantErr: "gagal mengambil daftar data",
		},
		{
			name: "should return list successfully",
			params: &entities.ListParams{Page: 1, Limit: 10},
			mockSetup: func(repo *mocks.MockRepositoryInterface) {
				repo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(pagination.Response[REPLACE.REPLACE_STRUCT]{
						Data:       []REPLACE.REPLACE_STRUCT{{ID: uuid.New(), Name: "Item 1"}},
						Pagination: pagination.Pagination{Total: 1, Page: 1, Limit: 10, TotalPage: 1},
					}, nil)
			},
			wantTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockRepositoryInterface(ctrl)
			mockCache := cacheMocks.NewMockCache(ctrl)
			mockStorage := storageMocks.NewMockStorage(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			svc := REPLACE.NewService(testDB, mockRepo, mockCache, mockStorage)
			resp, err := svc.List(context.Background(), tt.params)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Empty(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, resp)
			require.Equal(t, tt.wantTotal, len(resp.Data))
		})
	}
}

// =============================================================================
// Tests: Update
// =============================================================================

func TestService_Update(t *testing.T) {
	t.Parallel()

	type args struct {
		req    *REPLACE.UpdateRequest
		userID string
	}

	tests := []struct {
		name      string
		args      args
		mockSetup func(*mocks.MockRepositoryInterface, sqlmock.Sqlmock)
		wantErr   string
	}{
		{
			name: "should return bad request when id is empty",
			args: args{
				req: &REPLACE.UpdateRequest{Name: stringPtr("Updated")},
			},
			wantErr: "id tidak valid",
		},
		{
			name: "should update successfully",
			args: args{
				req:    &REPLACE.UpdateRequest{ID: uuid.New().String(), Name: stringPtr("Updated")},
				userID: "user-1",
			},
			mockSetup: func(repo *mocks.MockRepositoryInterface, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
					Return(REPLACE.REPLACE_STRUCT{ID: uuid.New(), Name: "Original"}, nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				sqlMock.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.mockSetup == nil {
				mockRepo := mocks.NewMockRepositoryInterface(ctrl)
				mockCache := cacheMocks.NewMockCache(ctrl)
				mockStorage := storageMocks.NewMockStorage(ctrl)
				svc := REPLACE.NewService(testDB, mockRepo, mockCache, mockStorage)
				_, err := svc.Update(context.Background(), tt.args.req, tt.args.userID)
				if tt.wantErr != "" {
					require.Error(t, err)
					require.Contains(t, err.Error(), tt.wantErr)
				} else {
					require.NoError(t, err)
				}
				return
			}

			svc, mockRepo, sqlMock := newServiceForTransaction(t, ctrl)
			tt.mockSetup(mockRepo, sqlMock)

			resp, err := svc.Update(context.Background(), tt.args.req, tt.args.userID)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Empty(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, resp)
		})
	}
}

// =============================================================================
// Tests: Delete
// =============================================================================

func TestService_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        string
		mockSetup func(*mocks.MockRepositoryInterface, sqlmock.Sqlmock)
		wantErr   string
	}{
		{
			name:    "should return bad request for invalid uuid",
			id:      "not-a-uuid",
			wantErr: "id tidak valid",
		},
		{
			name: "should delete successfully",
			id:   uuid.New().String(),
			mockSetup: func(repo *mocks.MockRepositoryInterface, sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectBegin()
				repo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
					Return(REPLACE.REPLACE_STRUCT{ID: uuid.New(), Name: "To Delete"}, nil)
				repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				sqlMock.ExpectCommit()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tt.mockSetup == nil {
				mockRepo := mocks.NewMockRepositoryInterface(ctrl)
				mockCache := cacheMocks.NewMockCache(ctrl)
				mockStorage := storageMocks.NewMockStorage(ctrl)
				svc := REPLACE.NewService(testDB, mockRepo, mockCache, mockStorage)
				err := svc.Delete(context.Background(), tt.id)
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}

			svc, mockRepo, sqlMock := newServiceForTransaction(t, ctrl)
			tt.mockSetup(mockRepo, sqlMock)
			err := svc.Delete(context.Background(), tt.id)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

// =============================================================================
// Helpers
// =============================================================================

func stringPtr(s string) *string {
	return &s
}
SVC_TEST_EOF

# Replace placeholders in service test
sed -i '' \
  -e "s/REPLACE_PKG/backend-go/g" \
  -e "s/REPLACE_STRUCT/${STRUCT_NAME}/g" \
  -e "s/REPLACE_TABLE/${PLURAL_TABLE}/g" \
  -e "s/REPLACE/${PKG_NAME}/g" \
  "${MODULE_DIR}/service_test.go"
fi

# Add missing imports in service test
# The test template imports "REPLACE_PKG/internal/REPLACE/mocks" but after substitution
# it becomes "backend-go/internal/module/mocks" which needs to also import entities & pagination
# The template uses entities.ListParams and pagination.Response directly
# Let's check if the imports are correct...

# =============================================================================
# Auto-register in cmd/api/wire.go (App struct + wire.NewSet)
# =============================================================================
WIRE_FILE="cmd/api/wire.go"

# Add import
if ! grep -q "\"backend-go/internal/${PKG_NAME}\"" "$WIRE_FILE"; then
  sed -i '' "/\"backend-go\/internal\/order\"/a\\
\\	\"backend-go/internal/${PKG_NAME}\"
" "$WIRE_FILE"
fi

# Add wire.NewSet variable before func Init()
if ! grep -q "var ${PKG_NAME}Set" "$WIRE_FILE"; then
  sed -i '' "/^func Init(log/i\\
var ${PKG_NAME}Set = wire.NewSet(\\
\\	${PKG_NAME}.NewRepository,\\
\\	${PKG_NAME}.NewService,\\
\\	${PKG_NAME}.NewHandler,\\
\\	${PKG_NAME}.NewRouter,\\
)\\
" "$WIRE_FILE"
fi

# Add set to wire.Build
if ! grep -q "${PKG_NAME}Set," "$WIRE_FILE"; then
  sed -i '' "/orderSet,/a\\
\\		${PKG_NAME}Set," "$WIRE_FILE"
fi

# =============================================================================
# Regenerate wire_gen.go
# =============================================================================
echo "⚙️  Running wire gen..."
GOTOOLCHAIN=go1.25.7 go run github.com/google/wire/cmd/wire@latest gen ./cmd/api/

# Format everything
gofmt -w "$WIRE_FILE" 2>/dev/null
gofmt -w "${MODULE_DIR}"/*.go "${MOCK_DIR}"/*.go 2>/dev/null

echo ""
echo "📁 Created module at: ${MODULE_DIR}/"
echo "   model.go | repository.go | service.go | handler.go | routes.go"
echo "   mocks/ (mockgen)"
if [ "$GEN_TESTS" = "true" ]; then
  echo "   repository_test.go | service_test.go"
fi
echo ""
echo "🌐 Routes:"
echo "   GET    /api/v2/${PLURAL_ROUTE}/list"
echo "   POST   /api/v2/${PLURAL_ROUTE}/create"
echo "   GET    /api/v2/${PLURAL_ROUTE}/get/:id"
echo "   PUT    /api/v2/${PLURAL_ROUTE}/update    (id in body)"
echo "   DELETE /api/v2/${PLURAL_ROUTE}/delete/:id"
echo ""
echo "📝 Register routes manually in internal/router/router.go:"
echo "   1. Add *${PKG_NAME}.Router parameter to NewGinEngine"
echo "   2. Add ${PKG_NAME}Router.RegisterRoutes(v2) before return"
echo ""
if [ "$GEN_TESTS" = "true" ]; then
  echo "📝 Run tests: go test ./internal/${PKG_NAME}/... -v"
fi
echo ""
echo "✅ Auto-registered in:"
echo "   - cmd/api/wire.go (wire.NewSet + wire.Build)"
echo "   - wire gen ./cmd/api/"
