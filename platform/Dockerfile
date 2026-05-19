FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -o /out/platform-api ./cmd/platform-api
RUN go build -o /out/platform-migrate ./cmd/platform-migrate
RUN go build -o /out/platform-seed ./cmd/platform-seed

FROM alpine:3.23 AS runtime
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/platform-api /usr/local/bin/platform-api
COPY --from=build /out/platform-migrate /usr/local/bin/platform-migrate
COPY --from=build /out/platform-seed /usr/local/bin/platform-seed
EXPOSE 8080
CMD ["platform-api"]
