📋 Goose Migration Rules
Panduan wajib untuk semua developer sebelum membuat atau mengubah file migration.

⚡ Golden Rule
Selalu git pull --rebase origin develop + goose up sebelum membuat migration baru.

🛑 Wajib Pakai IF NOT EXISTS / IF EXISTS
Semua statement migration harus idempotent — aman dijalankan berkali-kali tanpa error.
ADD COLUMN
-- ❌ SALAH
ALTER TABLE ms_users ADD COLUMN verified_at TIMESTAMP;

-- ✅ BENAR
ALTER TABLE ms_users ADD COLUMN IF NOT EXISTS verified_at TIMESTAMP;

DROP COLUMN
-- ❌ SALAH
ALTER TABLE ms_users DROP COLUMN verified_at;

-- ✅ BENAR
ALTER TABLE ms_users DROP COLUMN IF EXISTS verified_at;

CREATE TABLE
-- ❌ SALAH
CREATE TABLE ms_products (...);

-- ✅ BENAR
CREATE TABLE IF NOT EXISTS ms_products (...);

DROP TABLE
-- ❌ SALAH
DROP TABLE ms_products;

-- ✅ BENAR
DROP TABLE IF EXISTS ms_products;

CREATE INDEX
-- ❌ SALAH
CREATE INDEX idx_users_email ON ms_users(email);

-- ✅ BENAR
CREATE INDEX IF NOT EXISTS idx_users_email ON ms_users(email);


📁 Struktur File Migration
Setiap file migration wajib memiliki section Up dan Down:
-- +goose Up
ALTER TABLE ms_users
ADD COLUMN IF NOT EXISTS verified_at TIMESTAMP NULL,
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;

-- +goose Down
ALTER TABLE ms_users
DROP COLUMN IF EXISTS verified_at,
DROP COLUMN IF EXISTS is_active;


🔄 Flow Sebelum Membuat Migration
1. git pull
2. goose up
3. goose create <nama_migration> sql
4. Tulis SQL dengan IF NOT EXISTS / IF EXISTS
5. Test: goose down 1 → goose up
6. Commit & push


🚀 Flow Deploy ke Staging / Production
deploy branch baru → goose up

Jangan pernah goose down di staging/production.
 Kalau ada kesalahan, buat migration baru yang melakukan reverse.
Environment
Command
Catatan
Local dev
goose up / goose down
Bebas
Staging
goose up only
Jangan down
Production
goose up only
Jangan down


✅ Checklist Sebelum PR
[ ] Sudah git pull --rebase origin develop sebelum buat migration
[ ] Semua ADD COLUMN pakai IF NOT EXISTS
[ ] Semua DROP COLUMN pakai IF EXISTS
[ ] Semua CREATE TABLE pakai IF NOT EXISTS
[ ] Semua DROP TABLE pakai IF EXISTS
[ ] Semua CREATE INDEX pakai IF NOT EXISTS
[ ] Ada section -- +goose Down yang proper
[ ] Sudah test goose down 1 lalu goose up di local

❓ FAQ
Q: Kenapa harus IF NOT EXISTS?
 A: Karena state DB di tiap environment (local, staging, production) bisa tidak sinkron dengan goose_db_version. IF NOT EXISTS memastikan migration tidak crash meski kolom sudah ada.
Q: Boleh goose down di staging?
 A: Tidak. Kalau ada kesalahan, buat file migration baru yang melakukan reverse perubahan.
Q: Migration saya conflict dengan migration orang lain, gimana?
 A: git pull --rebase origin develop, rename file migration kamu ke timestamp terbaru, lalu goose down 1 + goose up di local.
Q: Jenkins gagal karena column already exists?
 A: Tambahkan IF NOT EXISTS di file migration tersebut, lalu push ulang.

