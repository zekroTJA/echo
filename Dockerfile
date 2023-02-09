FROM golang:1.20-alpine AS build
WORKDIR /build
COPY . .
RUN go build -o echo cmd/echo/main.go

FROM alpine:latest AS final
WORKDIR /app
COPY --from=build /build/echo echo
ENTRYPOINT ["/app/echo"]
