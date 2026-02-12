# 🐳 Docker for This Project - Simple Guide

> **TL;DR:** Docker is NOT required, but makes database setup easier. You can install PostgreSQL directly on your computer instead.

---

## 🤔 Do I Need Docker?

**Short Answer:** No, but it helps.

**Why Docker Helps:**
- Install PostgreSQL in 1 command (no complex setup)
- Same environment on any computer (Windows/Mac/Linux)
- Easy to reset/restart database
- Don't mess with your system

**Without Docker:**
- Need to install PostgreSQL manually
- Different process for Windows/Mac/Linux
- Need to manage PostgreSQL service manually
- But works perfectly fine!

---

## 📊 What We Use Docker For

In this project, Docker is ONLY used for:

### **1. PostgreSQL Database** (Main reason)
- Stores all our data (employees, documents)
- Required for the application

### **2. pgAdmin** (Optional - nice to have)
- Visual tool to view/manage database
- Like a GUI for database
- Can use command line instead

### **3. LocalStack** (Optional - for local AWS testing)
- Simulates AWS S3/SES locally
- Only for development
- Can connect to real AWS instead

---

## 🎯 Two Options: Choose One

### **Option A: Use Docker** (Recommended for beginners)

**Pros:**
- Quick setup (2 commands)
- Same setup for everyone
- Easy to reset if something breaks

**Cons:**
- Need to install Docker (one-time)
- Uses some computer memory

**Setup time:** 10 minutes

---

### **Option B: Install PostgreSQL Directly**

**Pros:**
- No Docker needed
- Uses less memory
- Traditional approach

**Cons:**
- Setup varies by OS
- Need to manage service
- Harder to reset

**Setup time:** 20-30 minutes

---

## 🚀 Option A: Docker Setup (Step by Step)

### **Step 1: Install Docker**

**Windows:**
1. Download Docker Desktop: https://www.docker.com/products/docker-desktop/
2. Run installer
3. Restart computer
4. Verify: Open Command Prompt, type `docker --version`

**Mac:**
1. Download Docker Desktop: https://www.docker.com/products/docker-desktop/
2. Drag to Applications folder
3. Open Docker Desktop
4. Verify: Open Terminal, type `docker --version`

**Linux (Ubuntu):**
```bash
sudo apt update
sudo apt install docker.io docker-compose
sudo systemctl start docker
sudo systemctl enable docker
docker --version
```

---

### **Step 2: Create Docker Configuration**

Create a file named `docker-compose.yml` in your project root:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: manpower-postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: manpower_dev
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

**What this means:**
- `postgres:15-alpine` - PostgreSQL version 15 (lightweight)
- `POSTGRES_USER/PASSWORD` - Login credentials
- `POSTGRES_DB` - Database name
- `5432:5432` - Port (connects your computer to container)
- `volumes` - Saves data even after restart

---

### **Step 3: Start PostgreSQL**

```bash
# Start database
docker-compose up -d

# Check if running
docker ps

# You should see:
# CONTAINER ID   IMAGE               STATUS          PORTS
# abc123...      postgres:15-alpine  Up 10 seconds   0.0.0.0:5432->5432
```

**What `-d` means:** Run in background (detached mode)

---

### **Step 4: Verify It's Working**

**Method 1: Using Docker command**
```bash
docker exec -it manpower-postgres psql -U postgres
# You should see: postgres=#
# Type: \l (to list databases)
# Type: \q (to quit)
```

**Method 2: Using psql (if installed)**
```bash
psql -U postgres -h localhost -d manpower_dev
# Password: postgres
```

---

### **Step 5: Common Docker Commands**

```bash
# Start database
docker-compose up -d

# Stop database
docker-compose down

# View logs
docker-compose logs -f postgres

# Restart database
docker-compose restart

# Stop and remove everything (CAUTION: deletes data)
docker-compose down -v

# Check status
docker ps
```

---

### **Step 6 (Optional): Add pgAdmin**

If you want a visual tool, add this to your `docker-compose.yml`:

```yaml
version: '3.8'

services:
  postgres:
    # ... (same as before)
  
  pgadmin:
    image: dpage/pgadmin4:latest
    container_name: manpower-pgadmin
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@admin.com
      PGADMIN_DEFAULT_PASSWORD: admin
    ports:
      - "5050:80"
    depends_on:
      - postgres

volumes:
  postgres_data:
```

**Then:**
1. Run `docker-compose up -d`
2. Open browser: http://localhost:5050
3. Login: admin@admin.com / admin
4. Add server:
   - Host: `postgres` (not localhost!)
   - Port: 5432
   - Username: postgres
   - Password: postgres

---

## 🛠️ Option B: Install PostgreSQL Directly

### **Windows**

1. Download installer: https://www.postgresql.org/download/windows/
2. Run installer
3. During setup:
   - Port: 5432
   - Password: (choose and remember)
   - Default database: postgres
4. Verify in Command Prompt:
   ```cmd
   psql --version
   ```

5. Create database:
   ```cmd
   psql -U postgres
   # Password: (your chosen password)
   CREATE DATABASE manpower_dev;
   \q
   ```

---

### **Mac**

**Method 1: Using Homebrew** (Recommended)
```bash
# Install Homebrew if not installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install PostgreSQL
brew install postgresql@15

# Start service
brew services start postgresql@15

# Verify
psql --version

# Create database
createdb manpower_dev
```

**Method 2: Using Postgres.app**
1. Download: https://postgresapp.com/
2. Drag to Applications
3. Open Postgres.app
4. Click "Initialize"

---

### **Linux (Ubuntu/Debian)**

```bash
# Update package list
sudo apt update

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib

# Start service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Verify
psql --version

# Switch to postgres user and create database
sudo -u postgres psql
CREATE DATABASE manpower_dev;
\q
```

---

## 🔧 Connecting Your App to Database

### **With Docker:**

Your `.env` file:
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/manpower_dev?sslmode=disable
```

### **With Local PostgreSQL:**

Your `.env` file:
```env
# Windows/Mac/Linux
DATABASE_URL=postgres://postgres:your_password@localhost:5432/manpower_dev?sslmode=disable
```

**Replace `your_password` with the password you set during installation**

---

## 🆘 Troubleshooting

### **Docker Issues**

**"Cannot connect to Docker daemon"**
```bash
# Windows/Mac: Make sure Docker Desktop is running
# Linux: 
sudo systemctl start docker
```

**"Port 5432 already in use"**
```bash
# PostgreSQL already installed locally
# Either:
# 1. Stop local PostgreSQL, or
# 2. Change Docker port in docker-compose.yml:
#    ports:
#      - "5433:5432"  # Use port 5433 instead
```

**"Container keeps restarting"**
```bash
# View logs to see error
docker logs manpower-postgres
```

---

### **PostgreSQL Issues**

**"psql: could not connect to server"**
```bash
# Check if PostgreSQL is running

# Windows:
# Open Services, look for "postgresql"

# Mac (with Homebrew):
brew services list

# Linux:
sudo systemctl status postgresql
```

**"Database does not exist"**
```bash
# Create it manually
psql -U postgres
CREATE DATABASE manpower_dev;
\q
```

**"Password authentication failed"**
- Check your password
- Check DATABASE_URL in .env file
- Make sure username is correct

---

## 📝 What You Should Know

### **Basic Concepts**

**Container:** Like a lightweight virtual machine
- Runs PostgreSQL
- Isolated from your computer
- Can be started/stopped easily

**Image:** Template for container
- `postgres:15-alpine` is the PostgreSQL image
- Downloaded from Docker Hub (once)

**Volume:** Persistent storage
- Saves database data
- Survives container restart/removal

**Port:** Connection point
- `5432` is PostgreSQL's default port
- Your app connects to this port

---

### **You DON'T Need to Know**

- Dockerfile syntax
- Docker networking
- Container orchestration
- Kubernetes
- Advanced Docker features

**You ONLY need:**
1. `docker-compose up -d` (start)
2. `docker-compose down` (stop)
3. `docker ps` (check status)

That's it!

---

## 🎯 Recommendation

**For This Project:**

1. **If you're new to programming:** Use Docker
   - Easier setup
   - Less can go wrong
   - Same commands for everyone

2. **If you're comfortable with databases:** Use local PostgreSQL
   - One less tool to learn
   - Faster performance
   - More control

**Either way works perfectly fine!**

---

## 📚 Resources

**Docker:**
- Official Docs: https://docs.docker.com/get-started/
- Docker Compose: https://docs.docker.com/compose/

**PostgreSQL:**
- Official Docs: https://www.postgresql.org/docs/
- Tutorial: https://www.postgresqltutorial.com/

---

## ✅ Quick Check

After setup (Docker OR local), verify:

```bash
# Can connect to database?
psql -U postgres -h localhost -d manpower_dev
# Should show: manpower_dev=#

# In psql, check:
\l                    # List databases (should see manpower_dev)
\dt                   # List tables (empty for now)
\q                    # Quit
```

If this works, you're ready to start development! 🎉

---

## 🤝 Getting Help

**Docker Issues:**
- Docker Desktop has built-in help
- Check logs: `docker logs manpower-postgres`

**PostgreSQL Issues:**
- Check service is running
- Verify port 5432 is free
- Check password is correct

**Still Stuck?**
- Check the main PROJECT_SETUP.md guide
- Search error message on Google
- Most errors are permission or port conflicts

---

**Remember:** Docker is just a tool to run PostgreSQL. It's not magic, just convenient! If you prefer installing PostgreSQL directly, that's completely fine and will work exactly the same way.
