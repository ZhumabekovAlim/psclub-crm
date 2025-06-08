FROM golang:1.23-alpine as build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/web/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/app .
COPY --from=build /app/db/migrations ./migrations
ENV GIN_MODE=release
EXPOSE 4000
CMD ["./app"]
