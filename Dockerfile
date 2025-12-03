FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go-mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api cmd/api/main.go

FROM alpine:3.19

RUN apk --no-cache tzdata

COPY --from=builder /app/api .

COPY .env .

EXPOSE 8080

CMD [ "/api" ]



