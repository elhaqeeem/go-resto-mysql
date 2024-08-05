# Gunakan image Golang resmi sebagai base image
FROM golang:1.21 AS builder

# Set working directory
WORKDIR /app

# Salin file go mod dan go sum
COPY go.mod go.sum ./

# Download dependensi
RUN go mod download

# Salin sisa kode
COPY . .

# Debugging: list files in /app
RUN ls -l /app

# Build aplikasi
RUN go build -o main .

# Gunakan image Alpine sebagai base image untuk aplikasi yang dibangun
FROM alpine:latest  

# Install library yang diperlukan
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Salin binary dari stage builder
COPY --from=builder /app/main .

# Tentukan perintah yang akan dijalankan saat container dijalankan
CMD ["./main"]