**NAMA**: KHOIRUL RIJAL WICAKSONO

**NIM**: 122140234
**Mata Kuliah**: Pemrograman Web Lanjut  
**Institusi**: Institut Teknologi Sumatera

# **Create Order Saga - gRPC Microservices dengan Go**
---

## 📌 **Deskripsi Proyek**

Proyek ini mengimplementasikan pola **Create Order Saga** menggunakan arsitektur **microservices** dengan bahasa pemrograman **Go** dan komunikasi antar layanan menggunakan **gRPC**. Tujuannya adalah untuk memfasilitasi transaksi terdistribusi yang aman, konsisten, dan tetap dapat diandalkan meskipun terjadi kegagalan pada salah satu layanan.  
Setiap kegagalan akan ditangani melalui **mekanisme kompensasi otomatis**, seperti pembatalan pesanan atau pengembalian dana.

---

## 🧩 **Arsitektur Sistem**

```text
+-------------+     +----------------+     +----------------+     +----------------+
| Orchestrator| --> | Order Service  | --> | Payment Service| --> | Shipping Service|
+-------------+     +----------------+     +----------------+     +----------------+
      |                        |                   |                      |
      |<-------- Mekanisme kompensasi (cancel/refund) jika gagal -------->|
```

Setiap service berjalan secara independen pada server gRPC-nya masing-masing:

- **Order Service** → port `50051`  
- **Payment Service** → port `50052`  
- **Shipping Service** → port `50053`  
- **Orchestrator** bertindak sebagai klien yang memanggil semua layanan di atas secara berurutan.

---

## 🔧 **Layanan (Services)**

### 🟦 OrderService
- `CreateOrder(userId)` – Membuat pesanan baru
- `CancelOrder(orderId)` – Membatalkan pesanan jika transaksi gagal

### 🟥 PaymentService
- `ProcessPayment(orderId)` – Memproses pembayaran untuk pesanan
- `RefundPayment(paymentId)` – Mengembalikan dana jika pengiriman gagal

### 🟧 ShippingService
- `StartShipping(orderId)` – Memulai proses pengiriman barang
- `CancelShipping(shippingId)` – Membatalkan pengiriman jika terjadi kendala

---

## ▶️ **Cara Menjalankan**

Buka terminal terpisah untuk setiap service:

```bash
# Jalankan Order Service
go run order/server/main.go

# Jalankan Payment Service
go run payment/server/main.go

# Jalankan Shipping Service
go run shipping/server/main.go

# Jalankan Orchestrator dari direktori root
go run orchestrator/main.go
```

---

## 🧪 **Skenario Pengujian**

### ✅ **Alur Berhasil (Success Flow)**
Semua layanan berjalan dengan lancar. Output:  
`>>> All steps completed successfully`

---

### ❌ **Gagal pada Payment**
- Order berhasil dibuat
- Proses pembayaran gagal
- Sistem memicu `CancelOrder` sebagai kompensasi  
Output:  
`>>> Payment failed. Order has been cancelled.`

---

### ❌ **Gagal pada Shipping**
- Order berhasil dibuat
- Pembayaran berhasil
- Pengiriman gagal
- Sistem menjalankan kompensasi: `RefundPayment` dan `CancelOrder`  
Output:  
`>>> Shipping failed. Payment refunded and order cancelled.`

🔧 Untuk melakukan simulasi kegagalan, ubah nilai:
```go
status := "FAILED" 
atau 
status := "CANCELLED"
```
pada service yang ingin diuji.

---

## 🛠️ **Teknologi yang Digunakan**
- **Go (Golang)**: Bahasa pemrograman utama
- **gRPC**: Komunikasi efisien antar layanan
- **Protocol Buffers (proto3)**: Serialisasi data dan definisi layanan
- **Go Modules**: Manajemen dependensi
- **protoc**: Compiler untuk menghasilkan kode stub dari file `.proto`

---

## 📁 **Struktur Folder**

```text
create_order_saga/
├── go.mod
├── proto/              # File .proto dan hasil generate pb.go
├── order/server/       # Source code OrderService
├── payment/server/     # Source code PaymentService
├── shipping/server/    # Source code ShippingService
└── orchestrator/       # Source code Orchestrator (client)
```
