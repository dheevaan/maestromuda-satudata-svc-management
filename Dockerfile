# Stage 1: Build
FROM golang:1.21-alpine AS build

RUN apk --no-cache add tzdata git

WORKDIR /app

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    export PATH=$PATH:$(go env GOPATH)/bin

# Copy semua source code ke dalam container
COPY source/. ./

# Generate Swagger documentation
RUN swag init --parseDependency --parseInternal

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

# Stage 2: Final Image
FROM scratch AS finals

WORKDIR /app

# Copy timezone data
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary hasil build
COPY --from=build /app/app /app/app

# Copy file Swagger docs
COPY --from=build /app/docs /app/docs

# Copy file konfigurasi (pastikan `.env` ada di host sebelum build)
COPY source/.env /app/.env

# Set timezone
ENV TZ=Asia/Jakarta

# Expose port
EXPOSE 8900

# Jalankan aplikasi
CMD ["/app/app"]
