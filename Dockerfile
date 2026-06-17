# ==========================================
# STAGE 1: Compilation Environment
# ==========================================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies separately so docker builds stay lightning fast
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Compile a static, production-optimized binary file
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ==========================================
# STAGE 2: Minimal Execution Runtime
# ==========================================
FROM alpine:latest

WORKDIR /app

# Copy the lightweight executable binary from the builder stage
COPY --from=builder /app/main .

# Copy UI templates directory explicitly
COPY --from=builder /app/templates ./templates

# Document that our app container listens out on port 8000
EXPOSE 8000

# Fire up the engine
CMD ["./main"]