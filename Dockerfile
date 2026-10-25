FROM golang:1.27.3-alpine AS builder
RUN apk add --no-cache gcc g++ git
# Google's module storage intermittently returns 403 to some hosting IPs (e.g. Hetzner).
# "|" falls through to the next source on any error; go.sum still verifies every module.
ARG GOPROXY="https://proxy.golang.org|https://goproxy.cn|direct"
ENV GOPROXY=$GOPROXY
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd/server
RUN CGO_ENABLED=1 GOOS=linux go build -o bloggy

FROM alpine:3.24

# Create a non-root
RUN adduser -D noroot && mkdir -p /etc/sudoers.d \
        && echo "noroot ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/noroot \
        && chmod 0440 /etc/sudoers.d/noroot \
        && chown -R noroot:noroot /home/noroot
USER noroot

WORKDIR /app
COPY --from=builder /app/cmd/server/bloggy /app/bloggy
CMD ["/app/bloggy"]
