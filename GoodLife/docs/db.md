steps to install PostgreSQL, start the service immediately, and create the `goodlife` db.

## 1. Install PostgreSQL
Update your package list and install PostgreSQL along with its contrib packages. On modern Ubuntu versions (24.04), this installs PostgreSQL 16 by default.

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib -y
```

The service usually starts automatically upon installation, but you can ensure it is running immediately with:

```bash
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

Verify the status to confirm it is **active (running)**:
```bash
sudo systemctl status postgresql
```

To create a new database named `goodlife` with a dedicated user also named `goodlife`, follow these steps. This ensures your Go application connects securely without using the default `postgres` superuser.

### 1. Access the PostgreSQL Shell
Switch to the `postgres` system user and open the `psql` prompt:
```bash
sudo -i -u postgres
psql -d postgres
```

### 2. Create the User and Database
Run the following SQL commands inside the `psql` prompt. You will be prompted to enter a password for the new user.

```sql
-- Create the user 'goodlife' with a password prompt
CREATE USER goodlife WITH PASSWORD 'goodlife';

-- Create the database 'goodlife' owned by the user 'goodlife'
CREATE DATABASE goodlife OWNER goodlife;

-- Grant all privileges on the database to the user (redundant if OWNER is set, but explicit)
GRANT ALL PRIVILEGES ON DATABASE goodlife TO goodlife;
```

### 3. Grant Schema Permissions (PostgreSQL 15+)
If you are using PostgreSQL 15 or newer, the `public` schema no longer grants `CREATE` permission by default. You must explicitly grant it so your application can create tables:

```sql
-- Connect to the new database to set schema permissions
\c goodlife

-- Grant usage and creation rights on the public schema
GRANT ALL ON SCHEMA public TO goodlife;

-- Ensure future tables are accessible by the user
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO goodlife;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO goodlife;
```

Type `\q` to exit the prompt, then `exit` to leave the postgres user session.

### 4. Connect from Go
Update your Go application's connection string to use the new credentials:
`postgres://goodlife:goodlife@localhost/goodlife?sslmode=disable`


