# Deployment and CI/CD Guide

This guide covers production deployment and a simple GitHub Actions CI/CD flow for the AWS Admin panel.

## Prerequisites

- Linux server (Ubuntu 20.04+ recommended)
- Docker Engine + Docker Compose v2
- A public domain name (optional but recommended)
- SMTP credentials for email sending

## Production Deployment (Docker Compose)

### 1) Server setup

```bash
sudo apt update
sudo apt install -y docker.io docker-compose-plugin
sudo usermod -aG docker $USER
```

Log out and back in to apply group changes.

### 2) Clone the repo

```bash
git clone https://github.com/ratons127/Cloud-admin-Panel-latest.git
cd Cloud-admin-Panel-latest
```

### 3) Configure environment variables

Copy `.env.example` to `.env` and update values:

```bash
cp .env.example .env
```

Required variables:

- `DB_URL`
- `JWT_SECRET`
- `APP_BASE_URL`
- `SUPER_ADMIN_EMAILS`
- `SMTP_HOST`
- `SMTP_PORT`
- `SMTP_USER`
- `SMTP_PASS`
- `SMTP_FROM`

Set `APP_ENV=production`.

### 4) Build and start

```bash
docker compose up -d --build
```

The app should be available on port `6006` by default (see `docker-compose.yml`).

### 5) Verify

- Visit `http://YOUR_SERVER:6006`
- Log in with a super admin user
- Create a user and verify email delivery

## Optional: Reverse Proxy (Nginx)

If you want HTTPS and a clean domain:

1) Point your domain to the server.
2) Install Nginx.
3) Use a basic reverse proxy config:

```nginx
server {
  listen 80;
  server_name your-domain.com;

  location / {
    proxy_pass http://127.0.0.1:6006;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }
}
```

Then use Certbot for HTTPS.

## CI/CD (GitHub Actions + SSH)

This is a minimal flow that builds on the server using Docker Compose.

### 1) Create SSH keys

Generate an SSH key on your local machine:

```bash
ssh-keygen -t ed25519 -C "github-actions"
```

Add the public key to the server:

```bash
ssh-copy-id user@your-server
```

### 2) Add GitHub Secrets

In your GitHub repo settings, add:

- `SSH_HOST`
- `SSH_USER`
- `SSH_PRIVATE_KEY` (contents of your private key)

### 3) Example workflow (add manually)

Create `.github/workflows/deploy.yml` with:

```yaml
name: Deploy
on:
  push:
    branches: [ "main" ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Deploy over SSH
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.SSH_HOST }}
          username: ${{ secrets.SSH_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd /opt/Cloud-admin-Panel-latest
            git pull origin main
            docker compose up -d --build
```

### 4) Server layout

On the server, keep the repo in a stable path such as:

```
/opt/Cloud-admin-Panel-latest
```

Make sure `.env` exists on the server.

## Notes

- Do not commit `.env` to public repos unless you fully trust the audience.
- Use a strong `JWT_SECRET`.
- If email fails, check SMTP settings and port access on your server firewall.
