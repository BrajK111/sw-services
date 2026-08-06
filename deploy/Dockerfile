FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/sw-services ./cmd/sw-services

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/sw-services ./sw-services
COPY configs ./configs
EXPOSE 8091
ENTRYPOINT ["./sw-services"]
