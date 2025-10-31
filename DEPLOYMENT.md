# BondiHub Deployment Guide

This guide covers deploying the BondiHub application to production environments.

## 🚀 Quick Start with Docker

### Prerequisites
- Docker and Docker Compose installed
- Domain name configured (optional)
- SSL certificates (for production)

### 1. Clone and Setup
```bash
git clone <repository-url>
cd BondiHub
```

### 2. Environment Configuration
```bash
# Copy environment files
cp backend/env.example backend/.env
cp frontend/src/environments/environment.ts frontend/src/environments/environment.prod.ts

# Edit configuration files
nano backend/.env
nano frontend/src/environments/environment.prod.ts
```

### 3. Production Environment Variables
```bash
# backend/.env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=bondihub
DB_SSLMODE=require
JWT_SECRET=your-super-secure-jwt-secret-key
JWT_EXPIRES_IN=24h
PORT=8080
GIN_MODE=release
CLOUDINARY_CLOUD_NAME=your-cloud-name
CLOUDINARY_API_KEY=your-api-key
CLOUDINARY_API_SECRET=your-api-secret
MTN_MOMO_API_URL=https://api.momodeveloper.mtn.com
MTN_MOMO_API_KEY=your-production-mtn-key
MTN_MOMO_SUBSCRIPTION_KEY=your-production-subscription-key
AIRTEL_MONEY_API_URL=https://openapi.airtel.africa
AIRTEL_MONEY_CLIENT_ID=your-production-client-id
AIRTEL_MONEY_CLIENT_SECRET=your-production-client-secret
COMMISSION_RATE=0.05
FEATURED_LISTING_PRICE=500.00
```

### 4. Build and Deploy
```bash
# Build frontend
cd frontend
npm install
npm run build --prod
cd ..

# Start services
docker-compose up -d
```

### 5. Verify Deployment
```bash
# Check services
docker-compose ps

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend
```

## 🌐 Manual Deployment

### Backend Deployment (Go)

#### 1. Server Setup
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database
sudo -u postgres psql
CREATE DATABASE bondihub;
CREATE USER bondihub_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE bondihub TO bondihub_user;
\q
```

#### 2. Deploy Application
```bash
# Clone repository
git clone <repository-url>
cd BondiHub/backend

# Install dependencies
go mod download

# Build application
go build -o bondihub main.go

# Create systemd service
sudo nano /etc/systemd/system/bondihub.service
```

#### 3. Systemd Service Configuration
```ini
[Unit]
Description=BondiHub API Server
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/bondihub
ExecStart=/opt/bondihub/bondihub
Restart=always
RestartSec=5
Environment=GIN_MODE=release
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
```

#### 4. Start Service
```bash
sudo systemctl daemon-reload
sudo systemctl enable bondihub
sudo systemctl start bondihub
sudo systemctl status bondihub
```

### Frontend Deployment (Angular)

#### 1. Build for Production
```bash
cd frontend
npm install
npm run build --prod
```

#### 2. Deploy to Nginx
```bash
# Install Nginx
sudo apt install nginx -y

# Copy build files
sudo cp -r dist/bondihub-frontend/* /var/www/html/

# Configure Nginx
sudo nano /etc/nginx/sites-available/bondihub
```

#### 3. Nginx Configuration
```nginx
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;
    root /var/www/html;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # API proxy
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Angular routes
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Static assets caching
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

#### 4. Enable Site
```bash
sudo ln -s /etc/nginx/sites-available/bondihub /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 🔒 SSL Configuration

### Using Let's Encrypt
```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain SSL certificate
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

## 📊 Monitoring and Logging

### 1. Application Logs
```bash
# View logs
sudo journalctl -u bondihub -f

# Log rotation
sudo nano /etc/logrotate.d/bondihub
```

### 2. Database Monitoring
```bash
# PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-*.log

# Database performance
sudo -u postgres psql
SELECT * FROM pg_stat_activity;
```

### 3. System Monitoring
```bash
# Install monitoring tools
sudo apt install htop iotop nethogs -y

# Monitor resources
htop
iotop
nethogs
```

## 🔧 Maintenance

### 1. Database Backups
```bash
# Create backup script
sudo nano /opt/backup-db.sh
```

```bash
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/opt/backups"
DB_NAME="bondihub"

mkdir -p $BACKUP_DIR
pg_dump -h localhost -U bondihub_user $DB_NAME > $BACKUP_DIR/bondihub_$DATE.sql
gzip $BACKUP_DIR/bondihub_$DATE.sql

# Keep only last 7 days
find $BACKUP_DIR -name "bondihub_*.sql.gz" -mtime +7 -delete
```

### 2. Application Updates
```bash
# Update application
cd /opt/bondihub
git pull origin main
go build -o bondihub main.go
sudo systemctl restart bondihub

# Update frontend
cd /opt/bondihub/frontend
git pull origin main
npm install
npm run build --prod
sudo cp -r dist/bondihub-frontend/* /var/www/html/
```

## 🚨 Troubleshooting

### Common Issues

#### 1. Database Connection Issues
```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check connection
psql -h localhost -U bondihub_user -d bondihub
```

#### 2. Application Won't Start
```bash
# Check logs
sudo journalctl -u bondihub -n 50

# Check configuration
./bondihub --help
```

#### 3. Frontend Not Loading
```bash
# Check Nginx status
sudo systemctl status nginx

# Check configuration
sudo nginx -t

# Check file permissions
sudo chown -R www-data:www-data /var/www/html/
```

## 📈 Performance Optimization

### 1. Database Optimization
```sql
-- Create indexes
CREATE INDEX idx_houses_status ON houses(status);
CREATE INDEX idx_houses_landlord ON houses(landlord_id);
CREATE INDEX idx_payments_agreement ON payments(agreement_id);
CREATE INDEX idx_rental_agreements_tenant ON rental_agreements(tenant_id);

-- Analyze tables
ANALYZE;
```

### 2. Application Optimization
```bash
# Enable Go optimizations
export GOOS=linux
export GOARCH=amd64
go build -ldflags="-s -w" -o bondihub main.go

# Use reverse proxy caching
# Add to Nginx config:
location /api/ {
    proxy_cache api_cache;
    proxy_cache_valid 200 5m;
    proxy_pass http://localhost:8080;
}
```

### 3. Frontend Optimization
```bash
# Enable production optimizations
ng build --prod --aot --build-optimizer --vendor-chunk --common-chunk

# Enable service worker
ng add @angular/pwa
```

## 🔐 Security Checklist

- [ ] Change default passwords
- [ ] Enable firewall (ufw)
- [ ] Configure SSL/TLS
- [ ] Set up fail2ban
- [ ] Regular security updates
- [ ] Database access restrictions
- [ ] API rate limiting
- [ ] Input validation
- [ ] CORS configuration
- [ ] Security headers

## 📞 Support

For deployment issues:
- Check logs: `sudo journalctl -u bondihub -f`
- Database: `sudo -u postgres psql`
- Nginx: `sudo nginx -t && sudo systemctl status nginx`
- System: `htop` and `df -h`
