FROM golang:1.21 AS builder
ARG SERVICE_NAME
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY ./lib ./lib
COPY ./services ./services
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/service ./services/$SERVICE_NAME

FROM alpine:3.18
COPY --from=builder /src/bin /bin
CMD ["service"]