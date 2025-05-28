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
- Nextjs

## How To Run Backend

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
CREATE USER <YOUR_USERNAME> WITH PASSWORD <YOUR_PASSWORD>;
GRANT ALL PRIVILEGES ON DATABASE qrify TO <YOUR_USERNAME>;
```

> **Note:** Replace `<YOUR_USERNAME>` and `<YOUR_PASSWORD>` with your actual PostgreSQL username and password.

- In the `backend` folder, create a `.env` file with the following content:

```
DB_USER=<YOUR_USERNAME>
DB_PASSWORD=<YOUR_PASSWORD>
```

4. run the project

- run the frontend:

```bash
cd frontend
npm install
npm run dev
```

- run the backend:

```bash
cd backend/cmd/api/
go run main.go
```

## How to Run the Frontend

1. Navigate to the frontend directory:
   ```sh
   cd frontend
   ```
2. Install dependencies:
   ```sh
   npm install
   ```
3. Start the development server:
   ```sh
   npm run dev
   ```
4. Open your browser and go to [http://localhost:3000](http://localhost:3000)

> **Note:** Make sure you have Node.js and npm installed. The default port is 3000. To use a different port, run:
>
> ```sh
> PORT=4000 npm run dev
> ```
>
> and visit [http://localhost:4000](http://localhost:4000)
