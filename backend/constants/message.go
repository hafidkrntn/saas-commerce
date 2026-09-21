package constants

// Error messages (user-facing, in Bahasa Indonesia)
const (
	ErrUnauthorized   = "tidak memiliki akses"
	ErrForbidden      = "akses ditolak"
	ErrNotFound       = "data tidak ditemukan"
	ErrInternalServer = "kesalahan pada server"
	ErrInvalidInput   = "format input tidak valid"
	ErrTokenExpired   = "token telah kedaluwarsa"
	ErrTokenInvalid   = "token tidak valid"
)

// Success messages
const (
	MsgDataFound  = "berhasil menampilkan data"
	MsgDataAdd    = "berhasil menambahkan data"
	MsgDataUpdate = "berhasil memperbarui data"
	MsgDataDelete = "berhasil menghapus data"
)
