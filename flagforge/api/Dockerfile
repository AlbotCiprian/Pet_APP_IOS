FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/bin/api ./bin/api
RUN mkdir -p /app/migrations
ENV API_PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/bin/api"]
