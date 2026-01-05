# Luồng Tạo Token (Authentication Flow)

## Tổng quan

Tài liệu này mô tả chi tiết luồng tạo và xác thực JWT token trong hệ thống todolist.

---

## 1. Luồng Tạo Token (Login)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           LUỒNG TẠO TOKEN (LOGIN)                               │
└─────────────────────────────────────────────────────────────────────────────────┘

┌──────────┐      POST /v1/login         ┌─────────────────────┐
│  Client  │ ──────────────────────────► │   Login Handler     │
│(Postman) │    {email, password}        │ (login_handler.go)  │
└──────────┘                             └──────────┬──────────┘
                                                    │
                                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│ BƯỚC 1: PARSE REQUEST                                                           │
│ ─────────────────────                                                           │
│ • ShouldBind(&loginUserData) - Parse JSON body thành struct UserLogin           │
│ • UserLogin chứa: email, password                                               │
└─────────────────────────────────────────────────────────────────────────────────┘
                                                    │
                                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│ BƯỚC 2: KHỞI TẠO DEPENDENCIES                                                   │
│ ─────────────────────────────                                                   │
│ • store = storage.NewSQLStore(db)     → Tạo storage để query DB                 │
│ • bcrypt = common.NewBcryptHasher(10) → Tạo hasher để verify password           │
│ • business = NewLoginBusiness(...)    → Tạo business logic layer                │
└─────────────────────────────────────────────────────────────────────────────────┘
                                                    │
                                                    ▼
                                         ┌──────────────────────┐
                                         │   LoginBusiness      │
                                         │   (login.go)         │
                                         └──────────┬───────────┘
                                                    │
                                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│ BƯỚC 3: TÌM USER TRONG DATABASE                                                 │
│ ───────────────────────────────                                                 │
│ • FindUser(ctx, {email: data.Email})                                            │
│ • Query: SELECT * FROM users WHERE email = ?                                    │
│ • Trả về: *model.User (chứa Id, Email, Password, Salt, Role, Status...)         │
│                                                                                 │
│ Nếu không tìm thấy → Return ErrEmailOrPasswordInvalid                           │
└─────────────────────────────────────────────────────────────────────────────────┘
                                                    │
                                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│ BƯỚC 4: VERIFY PASSWORD                                                         │
│ ───────────────────────                                                         │
│ • hasher.Compare(user.Password, data.Password + user.Salt)                      │
│ • So sánh: password từ DB với (password nhập + salt) đã hash                    │
│                                                                                 │
│ Nếu không khớp → Return ErrEmailOrPasswordInvalid                               │
└─────────────────────────────────────────────────────────────────────────────────┘
                                                    │
                                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│ BƯỚC 5: TẠO TOKEN PAYLOAD                                                       │
│ ─────────────────────────                                                       │
│ payload := &common.TokenPayload{                                                │
│     UId:   user.Id,           → ID của user trong DB                            │
│     URole: user.Role.String() → Role của user (user/admin/mod/shipper)          │
│ }                                                                               │
└─────────────────────────────────────────────────────────────────────────────────┘
                                                    │
                                                    ▼
                                         ┌──────────────────────┐
                                         │   JWT Provider       │
                                         │   (jwt/jwt.go)       │
                                         └──────────┬───────────┘
                                                    │
                                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│ BƯỚC 6: GENERATE JWT TOKEN                                                      │
│ ──────────────────────────                                                      │
│ tokenProvider.Generate(payload, expiry=30days)                                  │
│                                                                                 │
│ JWT Claims Structure:                                                           │
│ ┌─────────────────────────────────────────────────────────────────────────────┐ │
│ │ {                                                                           │ │
│ │   "payload": {                                                              │ │
│ │     "user_id": 123,          ← UId từ TokenPayload                          │ │
│ │     "role": "user"           ← URole từ TokenPayload                        │ │
│ │   },                                                                        │ │
│ │   "exp": 1234567890,         ← Thời gian hết hạn (now + 30 ngày)            │ │
│ │   "iat": 1234567890,         ← Thời gian tạo token                          │ │
│ │   "jti": "1234567890"        ← JWT ID (UnixNano timestamp)                  │ │
│ │ }                                                                           │ │
│ └─────────────────────────────────────────────────────────────────────────────┘ │
│                                                                                 │
│ • Ký token bằng HMAC-SHA256 với SYSTEM_SECRET                                   │
│ • myToken = t.SignedString([]byte(secret))                                      │
└─────────────────────────────────────────────────────────────────────────────────┘
                                                    │
                                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│ BƯỚC 7: TRẢ VỀ RESPONSE                                                         │
│ ───────────────────────                                                         │
│ {                                                                               │
│   "data": {                                                                     │
│     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",                         │
│     "created": "2024-01-05T10:30:00Z",                                          │
│     "expiry": 2592000  (30 ngày tính bằng giây)                                 │
│   }                                                                             │
│ }                                                                               │
└─────────────────────────────────────────────────────────────────────────────────┘
                                                    │
                                                    ▼
                                              ┌──────────┐
                                              │  Client  │
                                              └──────────┘
```

---

## 2. Chi Tiết Tác Dụng Của Mỗi Thực Thể

| Thực thể | File | Tác dụng |
|----------|------|----------|
| **Login Handler** | `module/user/transport/gin/login_handler.go` | Nhận HTTP request, parse body, khởi tạo dependencies, gọi business logic |
| **UserLogin** | `module/user/model/user.go` | Struct chứa email/password từ request |
| **LoginBusiness** | `module/user/biz/login.go` | Xử lý logic nghiệp vụ: tìm user, verify password, gọi token provider |
| **LoginStorage** | `module/user/storage/get.go` | Interface để query user từ PostgreSQL |
| **Hasher (Bcrypt)** | `common/hasher.go` | So sánh password đã hash với password nhập vào |
| **TokenPayload** | `common/token.go` | Struct chứa thông tin cần đưa vào JWT (user_id, role) |
| **Provider Interface** | `component/tokenprovider/provider.go` | Interface định nghĩa Generate() và Validate() |
| **jwtProvider** | `component/tokenprovider/jwt/jwt.go` | Implementation cụ thể: tạo JWT với HS256, ký bằng secret |
| **myClaim** | `component/tokenprovider/jwt/jwt.go` | Struct kết hợp TokenPayload + RegisteredClaims (exp, iat, jti) |
| **Token** | `component/tokenprovider/jwt/jwt.go` | Struct trả về chứa token string, created time, expiry |

---

## 3. Luồng Validate Token (Khi Gọi API Cần Auth)

```
┌──────────┐  Authorization: Bearer <token>   ┌─────────────────────┐
│  Client  │ ───────────────────────────────► │  RequireAuth        │
└──────────┘                                  │  (middleware)       │
                                              └──────────┬──────────┘
                                                         │
                    ┌────────────────────────────────────┼────────────────────────────────────┐
                    ▼                                    ▼                                    ▼
           ┌─────────────────┐                 ┌─────────────────┐                  ┌─────────────────┐
           │ Extract Token   │                 │ Validate Token  │                  │ Find User by ID │
           │ from Header     │                 │ (JWT Parse)     │                  │ from Database   │
           └─────────────────┘                 └─────────────────┘                  └─────────────────┘
                    │                                    │                                    │
                    ▼                                    ▼                                    ▼
           Tách "Bearer xxx"                   Parse JWT, verify                    Query: WHERE id =
           lấy token string                   signature & expiry                   payload.UserId()
                                                         │                                    │
                                                         ▼                                    ▼
                                              Trả về TokenPayload                   Set user vào context
                                              {user_id, role}                       c.Set("current_user", user)
                                                                                              │
                                                                                              ▼
                                                                                    Gọi c.Next() → Handler
```

### Chi tiết các bước Validate:

1. **Extract Token**: Tách token từ header `Authorization: Bearer <token>`
2. **Validate Token**:
   - Parse JWT bằng secret key
   - Verify signature (HMAC-SHA256)
   - Kiểm tra thời gian hết hạn (exp)
3. **Find User**: Query database với user_id từ token payload
4. **Check Status**: Kiểm tra user có bị banned không (status == 0)
5. **Set Context**: Lưu user vào Gin context để handler sử dụng

---

## 4. Cấu Trúc JWT Token

### Header
```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

### Payload
```json
{
  "payload": {
    "user_id": 123,
    "role": "user"
  },
  "exp": 1234567890,
  "iat": 1234567890,
  "jti": "1234567890123456789"
}
```

### Signature
```
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  SYSTEM_SECRET
)
```

---

## 5. Các API Cần Authentication

| Method | Endpoint | Middleware |
|--------|----------|------------|
| POST | `/v1/items` | RequireAuth |
| PATCH | `/v1/items/:id` | RequireAuth |
| DELETE | `/v1/items/:id` | RequireAuth |
| GET | `/v1/profile` | RequireAuth |

---

## 6. Cấu Hình

| Biến môi trường | Mô tả |
|-----------------|-------|
| `SYSTEM_SECRET` | Secret key để ký JWT token |
| Token Expiry | 30 ngày (2592000 giây) |
| Hash Cost | Bcrypt cost = 10 |
