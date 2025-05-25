# QRify - QR Code Generator with Analytics

QRify is a service that generates QR codes for URLs and tracks their usage through analytics. It provides a simple API to create QR codes and monitor their usage in real-time.

## Features

- Generate QR codes for any URL
- Track QR code scans with detailed analytics
- Set expiration times for QR codes
- Prometheus metrics integration

## Prerequisites

- Go 1.21 or later
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

```bash
go run cmd/api/main.go
```
