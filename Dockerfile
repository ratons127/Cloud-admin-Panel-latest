<<<<<<< HEAD
# build environment
FROM node:18-alpine AS frontend
WORKDIR /app
COPY frontend/package*.json /app/
RUN npm install
COPY frontend/index.html /app/index.html
COPY frontend/vite.config.js /app/vite.config.js
COPY frontend/src /app/src
COPY frontend/public /app/public
RUN npm run build

FROM golang:alpine
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app
COPY --from=frontend /app/dist /app/dist
RUN go build -o main ./backend
RUN addgroup -S appuser && adduser -S -D -H -h /app -G appuser appuser
RUN mkdir -p /app/tmp && chown -R appuser:appuser /app/tmp
USER appuser
CMD ["./main"]
=======
  # build environment
  FROM node:18-alpine AS frontend
  WORKDIR /app
  COPY frontend/package*.json /app/
  RUN npm install
  COPY frontend/index.html /app/index.html
  COPY frontend/vite.config.js /app/vite.config.js
  COPY frontend/src /app/src
  COPY frontend/public /app/public
  RUN npm run build

  FROM golang:alpine
  WORKDIR /app
  COPY . /app
  COPY --from=frontend /app/dist /app/dist
  RUN go build -o main ./backend
  USER appuser
  CMD ["./main"]

  Then rebuild and redeploy:

  docker build -t cloud-admin-panel .
  docker rm -f cloud-admin
  docker run -d --name cloud-admin \
    --network cloud-admin-net \
    -p 6005:6005 \
    -e DB_URL="postgres://app:app_pass@cloud-admin-db:5432/app_db?sslmode=disable" \
    -e JWT_SECRET="your-32+char-secret" \
    -e APP_BASE_URL="http://182.252.82.115:6005/auth" \
    -e APP_ENV="production" \
    -e SUPER_ADMIN_EMAILS="youradmin@example.com" \
    cloud-admin-panel
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
