FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker ./cmd/worker

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/bin/worker ./bin/worker
ENV WORKER_PORT=8090
EXPOSE 8090
ENTRYPOINT ["/app/bin/worker"]
