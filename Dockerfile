
FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod verify

COPY . .

RUN go build -o go-gallery ./main.go

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/go-gallery .
COPY --from=builder /app/docs /root/docs
COPY --from=builder /app/scripts /root/scripts
COPY --from=builder /app/sql /root/sql

EXPOSE 3000

CMD ["./go-gallery"]
