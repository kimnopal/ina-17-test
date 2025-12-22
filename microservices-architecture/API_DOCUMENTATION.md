# API Documentation

Dokumentasi lengkap untuk API Microservices Architecture (INA 17)

## Base URLs

- **User Service**: `{{user_url}}`
- **Booking Service**: `{{booking_url}}`
- **Payment Service**: `{{payment_url}}`

## Authentication

API menggunakan Bearer Token Authentication. Token didapatkan setelah melakukan login dan harus disertakan di header:

```
Authorization: Bearer <token>
```

---

## User Service

### 1. Register User

Mendaftarkan user baru ke sistem.

**Endpoint:** `POST /users`

**Request Body:**

```json
{
  "username": "admin",
  "password": "password"
}
```

**Response Success (201):**

```json
{
  "id": "uuid",
  "username": "admin",
  "created_at": "timestamp"
}
```

---

### 2. Login

Melakukan autentikasi user dan mendapatkan token.

**Endpoint:** `POST /login`

**Request Body:**

```json
{
  "username": "admin",
  "password": "password"
}
```

**Response Success (200):**

```json
{
  "token": "bearer_token_string",
  "user": {
    "id": "uuid",
    "username": "admin"
  }
}
```

---

### 3. Get Auth User

Mendapatkan informasi user yang sedang login.

**Endpoint:** `GET /users/auth`

**Headers:**

```
Authorization: Bearer <token>
```

**Response Success (200):**

```json
{
  "id": "uuid",
  "username": "admin",
  "created_at": "timestamp"
}
```

---

## Booking Service

### 1. Get All Events

Mendapatkan daftar semua event yang tersedia.

**Endpoint:** `GET /events`

**Response Success (200):**

```json
[
  {
    "id": "uuid",
    "name": "Event Name",
    "description": "Event Description",
    "date": "timestamp",
    "location": "Location"
  }
]
```

---

### 2. Get All Tickets by Event

Mendapatkan daftar tiket yang tersedia untuk suatu event.

**Endpoint:** `GET /events/:uuid/tickets`

**Path Parameters:**

- `uuid` (string, required): Event ID

**Example:** `GET /events/0cf33d20-ed2b-40e6-a72c-00878ca92b75/tickets`

**Response Success (200):**

```json
[
  {
    "id": "uuid",
    "event_id": "uuid",
    "ticket_type": "VIP",
    "price": 150000,
    "quantity": 100,
    "available": 95
  }
]
```

---

### 3. Create Booking

Membuat booking tiket baru.

**Endpoint:** `POST /bookings`

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "event_id": "0cf33d20-ed2b-40e6-a72c-00878ca92b75",
  "ticket_id": "f2c5d8e7-31be-4a0d-9222-c783eff43c23",
  "quantity": 1
}
```

**Response Success (201):**

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "event_id": "uuid",
  "ticket_id": "uuid",
  "quantity": 1,
  "total_price": 150000,
  "status": "PENDING",
  "created_at": "timestamp"
}
```

---

### 4. Get All Bookings

Mendapatkan daftar semua booking milik user yang sedang login.

**Endpoint:** `GET /bookings`

**Headers:**

```
Authorization: Bearer <token>
```

**Response Success (200):**

```json
[
  {
    "id": "uuid",
    "user_id": "uuid",
    "event_id": "uuid",
    "ticket_id": "uuid",
    "quantity": 1,
    "total_price": 150000,
    "status": "PENDING",
    "created_at": "timestamp"
  }
]
```

---

### 5. Get Booking Status

Mendapatkan status booking tertentu.

**Endpoint:** `GET /bookings/:uuid/status`

**Headers:**

```
Authorization: Bearer <token>
```

**Path Parameters:**

- `uuid` (string, required): Booking ID

**Example:** `GET /bookings/c8061234-da74-4856-8a26-6e7fbceb838f/status`

**Response Success (200):**

```json
{
  "booking_id": "uuid",
  "status": "PENDING"
}
```

**Status Values:**

- `PENDING`: Booking menunggu pembayaran
- `PAID`: Booking sudah dibayar
- `CANCELLED`: Booking dibatalkan

---

## Payment Service

### 1. Create Payment

Membuat pembayaran untuk booking yang sudah dibuat.

**Endpoint:** `POST /payments`

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "booking_id": "725161df-a6b2-4ee7-9580-0c4d0a2db676",
  "amount": 150000,
  "payment_method": "QRIS"
}
```

**Payment Methods:**

- `QRIS`
- `CREDIT_CARD`
- `BANK_TRANSFER`
- `E_WALLET`

**Response Success (201):**

```json
{
  "id": "uuid",
  "booking_id": "uuid",
  "user_id": "uuid",
  "amount": 150000,
  "payment_method": "QRIS",
  "status": "PENDING",
  "payment_url": "https://payment-gateway.com/pay/xxx",
  "created_at": "timestamp"
}
```

---

### 2. Get All Payments

Mendapatkan daftar semua pembayaran milik user yang sedang login.

**Endpoint:** `GET /payments`

**Headers:**

```
Authorization: Bearer <token>
```

**Response Success (200):**

```json
[
  {
    "id": "uuid",
    "booking_id": "uuid",
    "user_id": "uuid",
    "amount": 150000,
    "payment_method": "QRIS",
    "status": "PENDING",
    "created_at": "timestamp"
  }
]
```

**Payment Status Values:**

- `PENDING`: Menunggu pembayaran
- `PAID`: Sudah dibayar
- `FAILED`: Pembayaran gagal
- `EXPIRED`: Pembayaran kadaluarsa

---

### 3. Webhook Payment Gateway

Endpoint webhook untuk menerima notifikasi dari payment gateway ketika status pembayaran berubah.

**Endpoint:** `POST /payments/webhook/payment-gateway`

**Request Body:**

```json
{
  "payment_id": "a9a498be-90f9-4d32-a68d-f29a628ee1ad",
  "status": "PAID"
}
```

**Response Success (200):**

```json
{
  "message": "Webhook processed successfully"
}
```

**Notes:**

- Endpoint ini dipanggil oleh payment gateway, bukan oleh client
- Ketika status berubah menjadi `PAID`, payment service akan mengirim notifikasi ke booking service untuk mengupdate status booking

---

## Error Responses

Semua endpoint dapat mengembalikan error response dengan format:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

### Common Error Codes

- **400 Bad Request**: Request tidak valid atau data yang dikirim salah
- **401 Unauthorized**: Token tidak valid atau tidak ada
- **403 Forbidden**: User tidak memiliki akses ke resource
- **404 Not Found**: Resource tidak ditemukan
- **409 Conflict**: Konflik data (misalnya username sudah digunakan)
- **500 Internal Server Error**: Kesalahan server

---

## Flow Diagram

### Booking Flow

1. User melakukan login → mendapat token
2. User melihat daftar event → `GET /events`
3. User melihat tiket untuk event tertentu → `GET /events/:uuid/tickets`
4. User membuat booking → `POST /bookings` (status: PENDING)
5. User membuat payment → `POST /payments` (payment status: PENDING)
6. Payment gateway memproses pembayaran
7. Payment gateway mengirim webhook → `POST /payments/webhook/payment-gateway`
8. Payment service update status payment → PAID
9. Payment service notifikasi ke booking service
10. Booking service update status booking → PAID
11. User dapat mengecek status booking → `GET /bookings/:uuid/status`

---

## Notes

- Semua ID menggunakan format UUID
- Semua timestamp menggunakan format ISO 8601
- Currency dalam Rupiah (IDR)
- Semua endpoint yang memerlukan autentikasi harus menyertakan Bearer Token di header
