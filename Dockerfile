FROM golang:1.24-alpine AS builder

WORKDIR /myblog

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /myblog/myblog .

FROM alpine:3.20

WORKDIR /myblog

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /myblog/myblog ./myblog
COPY --from=builder /myblog/config ./config
COPY --from=builder /myblog/views ./views

RUN mkdir -p /myblog/log

EXPOSE 5678

CMD ["./myblog"]
