# BondiHub - House Renting Service for Zambia

A comprehensive house renting platform designed specifically for the Zambian market, featuring mobile money payments, local currency support, and mobile-first design.

## 🏗️ Architecture

### Backend (Go)
- **Framework**: Gin
- **ORM**: GORM + SQLC
- **Database**: PostgreSQL
- **Authentication**: JWT
- **File Storage**: Cloudinary
- **Payment**: MTN MoMo & Airtel Money APIs

### Frontend (Angular)
- **Framework**: Angular 17+
- **UI Library**: PrimeNG
- **Styling**: Tailwind CSS
- **State Management**: Angular Services + RxJS

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Docker (optional)

### Backend Setup
```bash
cd backend
go mod tidy
go run main.go
```

### Frontend Setup
```bash
cd frontend
npm install
ng serve
```

### Database Setup
```bash
# Using Docker
docker-compose up -d postgres

# Or manually
createdb bondihub
psql bondihub < migrations/001_initial_schema.sql
```

## 📱 Features

- **Multi-role Support**: Landlords, Tenants, Agents, Admins
- **Property Management**: CRUD operations with image uploads
- **Payment Integration**: MTN MoMo and Airtel Money
- **Review System**: Tenant reviews and ratings
- **Maintenance Requests**: Issue tracking and resolution
- **Favorites**: Save preferred properties
- **Notifications**: Real-time updates
- **Admin Dashboard**: Analytics and management

## 💰 Monetization

- Featured listings with automatic expiration
- Commission tracking on successful rentals
- Subscription plans for landlords/agents
- Advertising space management

## 🌍 Zambian Market Focus

- ZMW currency support
- Mobile money integration
- Mobile-first responsive design
- Low data usage optimization
- Local payment methods

## 📄 License

MIT License - see LICENSE file for details
