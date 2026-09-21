# API Reference — SaaS Ecommerce (LAKU)

Base URL: `http://localhost:8000/api/v2`

Semua endpoint (kecuali login/refresh) butuh header:
```
Authorization: Bearer <access_token>
```

Envelope respons:
```json
{ "code": 200, "data": { ... }, "msg": "berhasil menampilkan data" }
```

## Auth

| Method | Path | Deskripsi |
|---|---|---|
| POST | `/auth/login` | `{email, password}` → `{access_token, refresh_token, user}` |
| POST | `/auth/refresh` | `{refresh_token}` → pasangan token baru (rotasi) |
| POST | `/auth/logout` | `{refresh_token}` → cabut sesi |
| GET | `/auth/me` | Profil user yang login |

Seed user: `platform@laku.id` / `admin@laku.id` / `staff@laku.id` — password `admin123`.

## Platform admin (`/platform`, khusus role platform-admin)

| Method | Path | Deskripsi |
|---|---|---|
| GET/POST | `/platform/tenants` | List / buat tenant |
| GET/PUT/DELETE | `/platform/tenants/:id` | Detail / ubah / hapus tenant |
| GET | `/platform/plans` | Daftar paket SaaS |

## Catalog

| Method | Path | Deskripsi |
|---|---|---|
| GET/POST | `/categories` | List (search, pagination) / buat |
| GET/PUT/DELETE | `/categories/:id` | Detail / ubah / hapus |
| GET/POST | `/products` | List (`search`, `status`, `category_id`, `page`, `limit`) / buat |
| GET/PUT/DELETE | `/products/:id` | Detail (dengan images) / ubah / hapus |

Create product:
```json
{ "sku": "LPT-001", "name": "Laptop Pro", "price": 15000000, "cost_price": 12000000,
  "stock": 10, "low_stock_threshold": 3, "status": "active",
  "category_id": "uuid", "tags": ["laptop"], "images": ["https://..."] }
```

## Customers

| Method | Path | Deskripsi |
|---|---|---|
| GET/POST | `/customers` | List (`search`, `status`) / buat |
| GET/PUT/DELETE | `/customers/:id` | Detail (dengan addresses) / ubah / hapus |
| POST | `/customers/:id/addresses` | Tambah alamat |
| GET | `/customers/:id/activities` | Timeline customer |

## Orders

| Method | Path | Deskripsi |
|---|---|---|
| GET/POST | `/orders` | List (`status`, `search`, `date_from`, `date_to`) / buat |
| GET | `/orders/:id` | Detail + items + timeline |
| PATCH | `/orders/:id/status` | `{status}` — `pending→confirmed→processing→shipped→delivered`, atau `cancelled`/`refunded` |

Create order:
```json
{ "customer_id": "uuid", "customer_name": "Budi", "customer_email": "budi@x.com",
  "items": [{"product_id": "uuid", "quantity": 2}],
  "payment_method_id": "uuid", "shipping_method": "JNE", "shipping_cost": 25000,
  "discount": 0, "shipping_address": {"line1": "...", "city": "Jakarta"},
  "notes": "..." }
```
Stock otomatis di-reserve; customer stats diperbarui.

## Payments

| Method | Path | Deskripsi |
|---|---|---|
| GET/POST | `/payments/methods` | List / buat metode pembayaran |
| PUT/DELETE | `/payments/methods/:id` | Ubah / hapus |
| POST | `/payments/charge` | `{order_id}` → buat payment + redirect URL (mock) |
| GET | `/payments/order/:orderId` | Status pembayaran sebuah order |
| POST | `/payments/webhook` | `{transaction_id, status}` — mock provider callback (`paid`/`failed`) |

## Inventory

| Method | Path | Deskripsi |
|---|---|---|
| GET/POST | `/inventory/warehouses` | List / buat gudang |
| PUT/DELETE | `/inventory/warehouses/:id` | Ubah / hapus |
| GET | `/inventory/stock-movements` | List (`product_id`) |
| POST | `/inventory/stock-movements` | `{product_id, type: in\|out\|adjustment\|return, quantity}` — menyesuaikan stok produk |
| GET/POST | `/inventory/shipments` | List / buat shipment |
| PUT | `/inventory/shipments/:id` | Ubah; status `delivered` → stok masuk + movement |

## Settings

| Method | Path | Deskripsi |
|---|---|---|
| GET/PUT | `/settings/store` | Pengaturan toko |
| GET/PUT | `/settings/company` | Pengaturan perusahaan (address JSONB) |
| GET/POST | `/settings/tax-rates` | List / buat pajak |
| PUT/DELETE | `/settings/tax-rates/:id` | Ubah / hapus |
| GET/POST | `/settings/shipping-zones` | List / buat zona ongkir |
| PUT/DELETE | `/settings/shipping-zones/:id` | Ubah / hapus |
| GET | `/settings/notifications` | List notifikasi |
| PUT | `/settings/notifications/:id` | Ubah preferensi (email/sms/in_app) |

## Reviews

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/reviews` | List (`product_id`, `published=true\|false`) |
| POST | `/reviews` | `{product_id, customer_id?, customer_name, rating, title, content}` |
| PATCH | `/reviews/:id/moderate` | `{is_published, is_verified}` — memperbarui rating produk |
| DELETE | `/reviews/:id` | Hapus |

## Analytics

| Method | Path | Deskripsi |
|---|---|---|
| GET | `/analytics/summary` | KPI + % perubahan vs periode sebelumnya (`from`, `to` RFC3339) |
| GET | `/analytics/revenue` | Bulanan (`months`, default 12) |
| GET | `/analytics/sales` | Harian (`days`, default 30) |
| GET | `/analytics/top-products` | Produk terlaris (`limit`, default 5) |
| GET | `/analytics/top-categories` | Kategori terlaris |

## Users & Roles

| Method | Path | Deskripsi |
|---|---|---|
| GET/POST | `/users` | List / buat user (`role_id` wajib) |
| GET/PUT/DELETE | `/users/:id` | Detail / ubah / hapus |
| GET/POST | `/roles` | List / buat role (`permissions` bitmask) |
| GET/PUT/DELETE | `/roles/:id` | Detail / ubah (non-sistem) / hapus |

## Catatan

- **Multi-tenancy**: semua data difilter `tenant_id` dari klaim JWT. Token platform admin (tanpa tenant) hanya bisa mengakses `/platform/*`.
- **Webhook mock**: di lingkungan lokal webhook berada dalam grup terautentikasi; untuk provider asli perlu mekanisme verifikasi signature.
- **Stok**: order dibuat → `stock` berkurang & `reserved` bertambah; `shipped` → reserved dilepas; `cancelled` → stok dikembalikan.
