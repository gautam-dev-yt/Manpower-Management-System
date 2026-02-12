# 🚀 Manpower Management System - Quick Start Guide

This guide will help you set up the development environment for the Manpower Management System.

---

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

### **Required:**
- **Go 1.22+**: [Download](https://go.dev/dl/)
- **Node.js 20+**: [Download](https://nodejs.org/)
- **PostgreSQL 15+**: [Download](https://www.postgresql.org/download/)
- **Git**: [Download](https://git-scm.com/)

### **Optional but Recommended:**
- **Docker & Docker Compose**: [Download](https://www.docker.com/)
- **VSCode/Cursor**: [Download](https://code.visualstudio.com/) or [Cursor](https://cursor.sh/)
- **Make**: Usually pre-installed on Mac/Linux, [Windows install](http://gnuwin32.sourceforge.net/packages/make.htm)

### **Verify Installation:**

```bash
go version        # Should be 1.22+
node --version    # Should be 20+
psql --version    # Should be 15+
docker --version  # Optional
```

---

## 🏗️ Project Structure

```
manpower-management-system/
├── backend/          # Go API
│   ├── cmd/
│   ├── internal/
│   ├── migrations/
│   └── ...
├── frontend/         # Next.js App
│   ├── src/
│   ├── public/
│   └── ...
└── docs/            # Documentation
```

---

## 🎯 Option 1: Quick Start (Docker)

**Best for:** Quick setup, don't want to install PostgreSQL locally

```bash
# 1. Clone repository
git clone <your-repo-url>
cd manpower-management-system

# 2. Start services (PostgreSQL, pgAdmin, LocalStack)
docker-compose up -d

# 3. Setup backend
cd backend
cp .env.example .env
# Edit .env with your settings
make install-tools
make migrate-up
make run

# 4. Setup frontend (in new terminal)
cd frontend
npm install
cp .env.local.example .env.local
# Edit .env.local with your settings
npm run dev

# 5. Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# pgAdmin: http://localhost:5050 (admin@admin.com / admin)
```

---

## 🛠️ Option 2: Manual Setup

**Best for:** Full control, already have PostgreSQL installed

### **Step 1: Database Setup**

```bash
# Create database
createdb manpower_dev

# Or using psql
psql -U postgres
CREATE DATABASE manpower_dev;
\q
```

### **Step 2: Backend Setup**

```bash
# 1. Navigate to backend
cd backend

# 2. Copy environment file
cp .env.example .env

# 3. Edit .env file
# Update DATABASE_URL, AWS credentials, etc.
nano .env  # or use your preferred editor

# 4. Install Go tools
make install-tools
# OR manually:
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install golang.org/x/tools/cmd/goimports@latest

# 5. Install dependencies
go mod download

# 6. Run database migrations
make migrate-up

# 7. (Optional) Seed sample data
make seed

# 8. Run the application
make run
# OR
go run cmd/api/main.go

# Application should now be running at http://localhost:8080
```

### **Step 3: Frontend Setup**

```bash
# 1. Navigate to frontend (in new terminal)
cd frontend

# 2. Install dependencies
npm install

# 3. Copy environment file
cp .env.local.example .env.local

# 4. Edit .env.local
# Update NEXT_PUBLIC_API_URL if needed
nano .env.local

# 5. Run development server
npm run dev

# Application should now be running at http://localhost:3000
```

---

## 🔧 Development Workflow

### **Backend Commands**

```bash
# Run application
make run

# Run with hot reload (install air first)
go install github.com/cosmtrek/air@latest
air

# Run tests
make test

# Run tests with coverage
make test-cover

# Format code
make format

# Lint code
make lint

# Create new migration
make migrate-create name=add_users_table

# Generate sqlc code (if using sqlc)
make sqlc-generate

# Build binary
make build
./bin/manpower-api
```

### **Frontend Commands**

```bash
# Development server
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Lint code
npm run lint

# Format code
npm run format

# Type check
npm run type-check

# Run tests (when added)
npm run test
```

---

## 🗄️ Database Management

### **Using pgAdmin (Docker)**

1. Open http://localhost:5050
2. Login: `admin@admin.com` / `admin`
3. Add server:
   - Host: `postgres` (if using Docker) or `localhost`
   - Port: `5432`
   - Username: `postgres`
   - Password: `postgres`
   - Database: `manpower_dev`

### **Using psql**

```bash
# Connect to database
psql -U postgres -d manpower_dev

# Common commands
\dt                 # List tables
\d employees        # Describe table
\q                  # Quit
```

### **Migrations**

```bash
# Create new migration
make migrate-create name=add_documents_table

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Check migration status
migrate -path migrations -database "postgres://..." version
```

---

## 🌐 AWS Services Setup

### **For Local Development (LocalStack)**

If using Docker Compose, LocalStack provides local AWS services:

```bash
# Create S3 bucket
aws --endpoint-url=http://localhost:4566 s3 mb s3://manpower-files

# Verify SES sender
aws --endpoint-url=http://localhost:4566 ses verify-email-identity \
    --email-address alerts@example.com
```

### **For Production (Real AWS)**

1. **S3 Bucket:**
   ```bash
   aws s3 mb s3://your-bucket-name
   aws s3api put-bucket-cors --bucket your-bucket-name --cors-configuration file://cors.json
   ```

2. **SES Setup:**
   ```bash
   aws ses verify-email-identity --email-address your-email@example.com
   ```

3. **SNS (for SMS - optional):**
   ```bash
   aws sns set-sms-attributes --attributes DefaultSMSType=Transactional
   ```

---

## 🧪 Testing

### **Backend Tests**

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/service/... -v

# Run with coverage
make test-cover

# Integration tests
go test ./internal/repository/... -tags=integration
```

### **Frontend Tests (when added)**

```bash
# Unit tests
npm run test

# E2E tests
npm run test:e2e
```

---

## 🐛 Troubleshooting

### **Database Connection Issues**

```bash
# Check if PostgreSQL is running
pg_isready

# Check connection
psql -U postgres -d manpower_dev -c "SELECT 1"

# Reset database
dropdb manpower_dev
createdb manpower_dev
make migrate-up
```

### **Port Already in Use**

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

### **Go Module Issues**

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
go mod tidy
```

### **Frontend Issues**

```bash
# Clear npm cache
npm cache clean --force

# Delete node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

---

## 📝 Environment Configuration

### **Backend (.env)**

Key variables to configure:

```env
# Database
DATABASE_URL=postgres://user:pass@localhost:5432/manpower_dev?sslmode=disable

# AWS (use LocalStack for local dev)
AWS_REGION=us-east-1
AWS_ENDPOINT=http://localhost:4566  # Remove for production
S3_BUCKET=manpower-files

# Email
SES_FROM_EMAIL=alerts@yourdomain.com

# Notifications
NOTIFICATION_EMAIL_ENABLED=true
EXPIRY_CHECK_TIME=09:00
```

### **Frontend (.env.local)**

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

---

## 🚢 Deployment

### **Backend Deployment (AWS EC2)**

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/api cmd/api/main.go

# Copy to server
scp bin/api user@server:/opt/manpower-api/

# Setup systemd service
sudo systemctl start manpower-api
```

### **Frontend Deployment (Vercel)**

```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
vercel

# Production deployment
vercel --prod
```

---

## 📚 Additional Resources

- **Development Guidelines:** See `DEVELOPMENT_GUIDELINES.md` for detailed coding standards
- **API Documentation:** See `docs/API.md` (generate with Swagger/OpenAPI)
- **Database Schema:** See `docs/DATABASE.md`
- **Architecture:** See `docs/ARCHITECTURE.md`

---

## 🤝 Getting Help

If you encounter issues:

1. Check this guide's troubleshooting section
2. Review the logs: `docker-compose logs -f` or check console output
3. Verify environment variables are correct
4. Ensure all prerequisites are installed correctly
5. Check PostgreSQL is running: `pg_isready`

---

## ✅ Verification Checklist

After setup, verify everything works:

- [ ] Backend API responds: `curl http://localhost:8080/api/health`
- [ ] Frontend loads: Open http://localhost:3000
- [ ] Database accessible: `psql -U postgres -d manpower_dev`
- [ ] Can create employee via API
- [ ] Can view dashboard in frontend
- [ ] File upload works (if configured S3/LocalStack)

---

**You're all set! 🎉**

Start developing by reading the `DEVELOPMENT_GUIDELINES.md` for coding standards and architecture patterns.
