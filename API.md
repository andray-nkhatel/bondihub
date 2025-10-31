# BondiHub API Documentation

## Base URL
- Development: `http://localhost:8080/api/v1`
- Production: `https://api.bondihub.com/api/v1`

## Authentication
All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Response Format
All API responses follow this format:
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "error": null
}
```

## Error Codes
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `422` - Validation Error
- `500` - Internal Server Error

---

## 🔐 Authentication Endpoints

### Register User
```http
POST /auth/register
```

**Request Body:**
```json
{
  "full_name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "phone": "+260977123456",
  "role": "tenant"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "uuid",
      "full_name": "John Doe",
      "email": "john@example.com",
      "phone": "+260977123456",
      "role": "tenant",
      "is_active": true,
      "is_verified": false,
      "created_at": "2024-01-01T00:00:00Z"
    },
    "token": "jwt-token"
  }
}
```

### Login User
```http
POST /auth/login
```

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

### Get Profile
```http
GET /auth/profile
```

### Update Profile
```http
PUT /auth/profile
```

**Request Body:**
```json
{
  "full_name": "John Smith",
  "phone": "+260977654321",
  "profile_image": "https://cloudinary.com/image.jpg"
}
```

### Change Password
```http
PUT /auth/change-password
```

**Request Body:**
```json
{
  "current_password": "oldpassword",
  "new_password": "newpassword123"
}
```

---

## 🏠 House Endpoints

### Get All Houses
```http
GET /houses?page=1&limit=10&house_type=apartment&status=available&min_rent=1000&max_rent=5000&bedrooms=2&bathrooms=1&featured=true&search=keyword
```

**Query Parameters:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10)
- `house_type` - apartment, house, studio, townhouse, commercial
- `status` - available, occupied, maintenance
- `min_rent` - Minimum rent amount
- `max_rent` - Maximum rent amount
- `bedrooms` - Minimum bedrooms
- `bathrooms` - Minimum bathrooms
- `featured` - Show only featured houses (true/false)
- `search` - Search in title, description, address

**Response:**
```json
{
  "success": true,
  "message": "Houses retrieved successfully",
  "data": {
    "houses": [
      {
        "id": "uuid",
        "title": "Beautiful 3BR Apartment",
        "description": "Spacious apartment in Lusaka",
        "address": "123 Independence Avenue, Lusaka",
        "monthly_rent": 3500.00,
        "status": "available",
        "house_type": "apartment",
        "bedrooms": 3,
        "bathrooms": 2,
        "area": 120.5,
        "is_featured": true,
        "landlord": {
          "id": "uuid",
          "full_name": "Jane Landlord",
          "phone": "+260977123456"
        },
        "images": [
          {
            "id": "uuid",
            "image_url": "https://cloudinary.com/image1.jpg",
            "is_primary": true
          }
        ],
        "average_rating": 4.5,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 50,
      "total_pages": 5
    }
  }
}
```

### Get Single House
```http
GET /houses/{id}
```

### Create House (Landlord/Admin)
```http
POST /houses
```

**Request Body:**
```json
{
  "title": "Beautiful 3BR Apartment",
  "description": "Spacious apartment in Lusaka with modern amenities",
  "address": "123 Independence Avenue, Lusaka",
  "monthly_rent": 3500.00,
  "house_type": "apartment",
  "latitude": -15.3875,
  "longitude": 28.3228,
  "bedrooms": 3,
  "bathrooms": 2,
  "area": 120.5,
  "is_featured": false
}
```

### Update House (Landlord/Admin)
```http
PUT /houses/{id}
```

### Delete House (Landlord/Admin)
```http
DELETE /houses/{id}
```

### Upload House Image (Landlord/Admin)
```http
POST /houses/{id}/images
```

**Request:** Multipart form data with `image` file

### Delete House Image (Landlord/Admin)
```http
DELETE /houses/images/{imageId}
```

---

## 💰 Payment Endpoints

### Process Payment
```http
POST /payments
```

**Request Body:**
```json
{
  "agreement_id": "uuid",
  "amount": 3500.00,
  "method": "MTN",
  "reference_no": "PAY_123456789"
}
```

**Payment Methods:**
- `MTN` - MTN MoMo
- `Airtel` - Airtel Money
- `Cash` - Cash payment
- `Bank` - Bank transfer

### Get Payments
```http
GET /payments?page=1&limit=10&status=completed&method=MTN
```

### Get Payment Details
```http
GET /payments/{id}
```

### Get Payment Statistics
```http
GET /payments/stats
```

**Response:**
```json
{
  "success": true,
  "message": "Payment statistics retrieved successfully",
  "data": {
    "total_payments": 150,
    "total_amount": 525000.00,
    "completed_payments": 140,
    "completed_amount": 490000.00,
    "pending_payments": 8,
    "failed_payments": 2,
    "payments_by_method": [
      {
        "method": "MTN",
        "count": 80,
        "amount": 280000.00
      },
      {
        "method": "Airtel",
        "count": 60,
        "amount": 210000.00
      }
    ]
  }
}
```

---

## 📋 Rental Agreement Endpoints

### Create Rental Agreement (Landlord/Admin)
```http
POST /rentals
```

**Request Body:**
```json
{
  "house_id": "uuid",
  "tenant_id": "uuid",
  "start_date": "2024-02-01",
  "end_date": "2025-01-31",
  "rent_amount": 3500.00,
  "deposit": 3500.00
}
```

### Get Rental Agreements
```http
GET /rentals?page=1&limit=10&status=active
```

### Get Rental Agreement Details
```http
GET /rentals/{id}
```

### Update Rental Agreement
```http
PUT /rentals/{id}
```

**Request Body:**
```json
{
  "status": "terminated"
}
```

### Terminate Rental Agreement
```http
PUT /rentals/{id}/terminate
```

---

## ⭐ Review Endpoints

### Create Review (Tenant)
```http
POST /reviews
```

**Request Body:**
```json
{
  "house_id": "uuid",
  "rating": 5,
  "comment": "Great apartment, highly recommended!"
}
```

### Get House Reviews
```http
GET /houses/{id}/reviews?page=1&limit=10
```

### Get User Reviews
```http
GET /reviews/my?page=1&limit=10
```

### Update Review
```http
PUT /reviews/{id}
```

### Delete Review
```http
DELETE /reviews/{id}
```

---

## 🔧 Maintenance Request Endpoints

### Create Maintenance Request (Tenant)
```http
POST /maintenance
```

**Request Body:**
```json
{
  "house_id": "uuid",
  "title": "Broken Water Heater",
  "description": "The water heater in the bathroom is not working properly",
  "priority": "high"
}
```

**Priority Levels:**
- `low` - Low priority
- `medium` - Medium priority
- `high` - High priority
- `urgent` - Urgent

### Get Maintenance Requests
```http
GET /maintenance?page=1&limit=10&status=pending&priority=high
```

### Get Maintenance Request Details
```http
GET /maintenance/{id}
```

### Update Maintenance Request (Landlord/Admin)
```http
PUT /maintenance/{id}
```

**Request Body:**
```json
{
  "status": "in_progress"
}
```

**Status Options:**
- `pending` - Pending
- `in_progress` - In Progress
- `resolved` - Resolved
- `cancelled` - Cancelled

### Get Maintenance Statistics
```http
GET /maintenance/stats
```

---

## ❤️ Favorite Endpoints

### Add to Favorites (Tenant)
```http
POST /favorites/{houseId}
```

### Remove from Favorites (Tenant)
```http
DELETE /favorites/{houseId}
```

### Get Favorites (Tenant)
```http
GET /favorites?page=1&limit=10
```

### Check if House is Favorite (Tenant)
```http
GET /favorites/{houseId}/check
```

---

## 🔔 Notification Endpoints

### Get Notifications
```http
GET /notifications?page=1&limit=20&unread_only=true&type=payment
```

### Get Notification Details
```http
GET /notifications/{id}
```

### Mark as Read
```http
PUT /notifications/{id}/read
```

### Mark All as Read
```http
PUT /notifications/read-all
```

### Delete Notification
```http
DELETE /notifications/{id}
```

### Get Notification Statistics
```http
GET /notifications/stats
```

---

## 👑 Admin Endpoints

### Get Dashboard Statistics
```http
GET /admin/dashboard
```

### Get All Users
```http
GET /admin/users?page=1&limit=20&role=landlord&search=john
```

### Update User Status
```http
PUT /admin/users/{id}/status
```

**Request Body:**
```json
{
  "is_active": false
}
```

### Get Reports
```http
GET /admin/reports?type=payments&start_date=2024-01-01&end_date=2024-01-31
```

**Report Types:**
- `payments` - Payment reports
- `houses` - Property reports
- `users` - User reports

---

## 📊 Data Models

### User
```json
{
  "id": "uuid",
  "full_name": "string",
  "email": "string",
  "phone": "string",
  "role": "landlord|tenant|agent|admin",
  "is_active": boolean,
  "is_verified": boolean,
  "profile_image": "string",
  "subscription_plan": "basic|premium|enterprise",
  "plan_expiry_date": "datetime",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### House
```json
{
  "id": "uuid",
  "landlord_id": "uuid",
  "title": "string",
  "description": "string",
  "address": "string",
  "monthly_rent": number,
  "status": "available|occupied|maintenance",
  "house_type": "apartment|house|studio|townhouse|commercial",
  "latitude": number,
  "longitude": number,
  "bedrooms": number,
  "bathrooms": number,
  "area": number,
  "is_featured": boolean,
  "featured_until": "datetime",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Payment
```json
{
  "id": "uuid",
  "agreement_id": "uuid",
  "amount": number,
  "payment_date": "datetime",
  "method": "MTN|Airtel|Cash|Bank",
  "reference_no": "string",
  "status": "pending|completed|failed|refunded",
  "commission": number,
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

---

## 🔒 Rate Limiting

- **General API**: 100 requests per minute per IP
- **Authentication**: 5 login attempts per minute per IP
- **File Upload**: 10 uploads per minute per user
- **Payment Processing**: 5 payments per minute per user

## 📝 Pagination

All list endpoints support pagination:
- `page` - Page number (starts from 1)
- `limit` - Items per page (max 100)
- Response includes pagination metadata

## 🔍 Filtering and Sorting

Most list endpoints support:
- **Filtering** by various fields
- **Search** in text fields
- **Sorting** by creation date (newest first by default)

## 📱 Mobile Optimization

- All endpoints are optimized for mobile usage
- Image uploads support mobile-friendly formats
- Response sizes are minimized for low-bandwidth connections
- Caching headers are set appropriately for mobile apps
