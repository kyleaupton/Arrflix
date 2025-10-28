FROM golang:1.24 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/snaggle-ops .

FROM debian:stable-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=build /out/snaggle-ops /usr/local/bin/snaggle-ops
ENTRYPOINT ["/usr/local/bin/snaggle-ops"]
