## Server Setup

#### Database Configuration

**Install PostgreSQL:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# Start and enable PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

**Configure Authentication:**

PostgreSQL defaults to `ident` authentication, but applications need `md5` authentication. Modify the authentication configuration:

```bash
# Find and configure PostgreSQL authentication
HBA_FILE=$(sudo -u postgres psql -t -c "SHOW hba_file;" | xargs)
sudo cp "$HBA_FILE" "$HBA_FILE.backup"

# Modify authentication methods to md5
sudo sed -i 's/local   all             all                                     peer/local   all             all                                     md5/' "$HBA_FILE"
sudo sed -i 's/host    all             all             127.0.0.1\/32            ident/host    all             all             127.0.0.1\/32            md5/' "$HBA_FILE"
sudo sed -i 's/host    all             all             ::1\/128                 ident/host    all             all             ::1\/128                 md5/' "$HBA_FILE"

# Restart PostgreSQL to apply changes
sudo systemctl restart postgresql
```

**Create Database and User:**
```bash
sudo -u postgres psql
```

In PostgreSQL command line:
```sql
CREATE USER c3user WITH PASSWORD 'your-database-password';
CREATE DATABASE c3_db OWNER c3user;
GRANT ALL PRIVILEGES ON DATABASE c3_db TO c3user;
\q
```

**Important**: Replace `'your-database-password'` with a secure password of your choice. You will need to use the same password in your `.env` file's `DB_PASSWORD` setting.

#### Configure Environment Variables

```bash
# Use the template for native installation
cp .env.example .env
```

Edit `.env` file with your database password and authentication settings:
- Set `DB_PASSWORD` to match the password you used when creating the database user
- Set `AUTH_USERNAME` and `AUTH_PASSWORD` for web interface login
- Set `BASE_PATH` to the URL path prefix if deploying under a subdirectory (e.g., `/c3` for `http://domain.com/c3/`), leave empty for root deployment

**Note:** Other configuration parameters not detailed here are also important for specific use cases. Please refer to the complete [configuration reference](#configuration-reference) section below for all available options.

#### Start Server

```bash
# Linux:
./server
# Windows:
.\server.exe
```

#### Access Web Interface

Open browser and navigate to: **http://localhost:3000**


### configuration-reference

| Variable | Default | Description |
| `BASE_PATH` | `""` | Base URL path prefix |
| `HOST` | `0.0.0.0` | Server host |
| `PORT` | `3000` | Server port |
| `ENV` | `` | Reserved |
| `AUTH_ENABLED` | `true` | Enable authentication |
| `AUTH_USERNAME` | `admin` | Login username |
| `AUTH_PASSWORD` | - | Login password |
| `SESSION_EXPIRE_HOURS` | `24` | Session expiration |
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_NAME` | `c3_db` | Database name |
| `DB_USER` | `c3user` | Database username |
| `DB_PASSWORD` | - | Database password |
| `DB_LOGGING` | `false` | Database query logging |
| `UPLOAD_DIR` | `uploads` | Upload directory |
| `LOG_DIR` | `logs` | Log file directory |
| `LOG_LEVEL` | `info` | Logging level |