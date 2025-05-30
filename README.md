# QRify - QR Code Generator with Analytics

QRify is a full-stack web application that allows users to generate QR codes for any URL and track their usage with analytics. It is designed for ease of use, flexibility, and modern deployment.

---

## Features

- **Generate QR codes** for any URL
- **Track QR code scans** with how many time scanned
- **Set expiration times** for QR codes

---

## Tech Stack

- **Frontend:** Next.js (React, TypeScript)
- **Backend:** Go (Gin framework, REST API)
- **Database:** PostgreSQL
- **Metrics:** Prometheus
- **Reverse Proxy:** Nginx (with optional HTTPS via Let's Encrypt)
- **Containerization:** Docker & Docker Compose

---

## User are able to

1. **User visits the QRify web app** on desktop or mobile.
2. **User enters a URL** and selects an expiration time (optional).
3. **QRify generates a unique QR code** for the URL and displays it to the user.
4. **User can download the QR code** or copy its unique ID.
5. **Anyone scanning the QR code** is redirected to the original URL (unless expired).
6. **User can check analytics** (scan count, expiration status) for any QR code by entering its ID.

---

## Note

- For local development and production, environment variables are used for configuration (see `.env.local` and `.env.production` in the respective folders).
- For production deployments, Docker Compose and Nginx are recommended. HTTPS can be enabled with Let's Encrypt and Certbot for your custom domain.
