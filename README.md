# QRify - QR Code Generator with Analytics

QRify is a service that generates QR codes for URLs and tracks their usage through analytics. It provides a simple API to create QR codes and monitor their usage in real-time.

## Features

- Generate QR codes for any URL
- Track QR code scans with detailed analytics
- Set expiration times for QR codes
- Prometheus metrics integration

## Prerequisites

- Go 1.23 or later
- PostgreSQL
- Docker (optional)
- Kubernetes (optional)

## Quick Start

1. Clone the repository:

```bash
git clone https://github.com/phucnguyen/qrify.git
cd qrify
```

2. Install dependencies:

```bash
go mod download
```

3. Set up the database:

- Install and start PostgreSQL if you haven't already.
- Create a new database and user for QRify:

```
psql -U postgres

-- In the psql shell:
CREATE DATABASE qrify;
CREATE USER qrify_user WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE qrify TO qrify_user;
```

- In the `backend` folder, create a `.env` file with the following content:

```
DB_USER=qrify_user
DB_PASSWORD=postgres
```

4. run the project

- run the frontend:

```bash
cd frontend
npm run dev
```

- run the backend:

```bash
cd backend/cmd/api/
go run main.go
```
