# ---------- BUILD STAGE ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY app ./app
WORKDIR /app/app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myshell

# ---------- RUN STAGE ----------
FROM alpine:latest

WORKDIR /root
COPY --from=builder /app/app/myshell .

ENV HISTFILE=/root/.myshell_history

CMD ["./myshell"]

