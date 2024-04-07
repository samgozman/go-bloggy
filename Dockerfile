FROM golang:1.22.2-alpine AS builder
RUN apk add --no-cache gcc g++
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/server
RUN CGO_ENABLED=1 GOOS=linux go build -o bloggy

FROM alpine:3.19

# Create a non-root
RUN adduser -D noroot && mkdir -p /etc/sudoers.d \
        && echo "noroot ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/noroot \
        && chmod 0440 /etc/sudoers.d/noroot \
        && chown -R noroot:noroot /home/noroot
USER noroot

WORKDIR /app
COPY --from=builder /app/cmd/server/bloggy /app/bloggy
CMD ["/app/bloggy"]
