# Third Party Integrations

Folder ini berisi implementasi koneksi ke API/service external (third-party).

## Struktur

```
third_party/
├── payment/        # Contoh: Midtrans, Xendit, Stripe
│   └── payment.go
├── shipping/       # Contoh: JNE, SiCepat, GoSend
│   └── shipping.go
├── notification/   # Contoh: Firebase, OneSignal, WhatsApp API
│   └── notification.go
└── ...
```

## Conventions

- Setiap integrasi punya folder sendiri (`third_party/<nama_service>/`)
- Define interface di file utama untuk mockability
- Gunakan `pkg/httpclient` untuk HTTP calls
- Request/response DTOs didefinisikan di file yang sama
- Error handling menggunakan `pkg/apperror`

## Contoh penggunaan di service

```go
import "backend-go/third_party/payment"

type Service struct {
    paymentClient payment.ClientInterface
}
```
