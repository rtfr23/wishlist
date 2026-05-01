FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /wishlist_api .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /wishlist_api .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY migrations ./migrations

RUN apk --no-cache add ca-certificates
EXPOSE 8080
CMD ["sh", "-c", "migrate -path ./migrations -database \"$CONN_STRING\" up && ./wishlist_api"]