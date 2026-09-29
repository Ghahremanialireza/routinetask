# ---------- Stage 1: Build ----------
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ---------- Stage 2: Run ----------
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server .

RUN mkdir -p /data
ENV DB_PATH=/data/routinetask.db

EXPOSE 8080

CMD ["./server"]
