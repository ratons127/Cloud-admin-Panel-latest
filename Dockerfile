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
