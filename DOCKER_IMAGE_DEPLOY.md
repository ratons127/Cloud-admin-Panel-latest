# Deploy Using Docker Hub Image

This note explains how to deploy the app using the published Docker image.

Image:

- `raton127/aws-admin_panel:latest`
- `raton127/aws-admin_panel:v1.0.2`
  
## Docker Install 

```
curl -fsSL https://get.docker.com/ | sh

```

## 1) Create a server directory

```bash
mkdir -p /opt/aws-admin
cd /opt/aws-admin
```

## 2) Create `.env`

Create a file named `.env` with these values:

```ini
APP_ENV=production
DB_URL=postgres://app:app_pass@db:5432/app_db?sslmode=disable
JWT_SECRET=your-32+char-secret
APP_BASE_URL=http://YOUR_SERVER:6006
SUPER_ADMIN_EMAILS=admin@example.com

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=you@example.com
SMTP_PASS=app-password
SMTP_FROM=you@example.com
```

## 3) Create `docker-compose.yml`

```yaml
  version: "3.8"

  services:
    db:
      image: postgres:16
      container_name: cloud-admin-db
      environment:
        POSTGRES_USER: app
        POSTGRES_PASSWORD: app_pass
        POSTGRES_DB: app_db
      volumes:
        - cloud_admin_pgdata:/var/lib/postgresql/data
      networks:
        - cloud-admin-net
      healthcheck:
        test: ["CMD-SHELL", "pg_isready -U app -d app_db"]
        interval: 5s
        timeout: 5s
        retries: 10

    app:
      image: raton127/aws-admin_panel:latest
      container_name: cloud-admin
      depends_on:
        db:
          condition: service_healthy
      ports:
        - "6006:6005"
      env_file:
        - .env
      networks:
    cloud_admin_pgdata:

  networks:
    cloud-admin-net:
```

## 4) Start the stack

```bash
docker compose up -d
```

## 5) Access the app

Open:

```
http://YOUR_SERVER:6006
```

## 6) Update to a new version

```bash
docker pull raton127/aws-admin_panel:latest
docker compose up -d
```
