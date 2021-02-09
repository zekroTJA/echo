FROM go:1.15-alpine AS build
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o echo cmd/echo/main.go

FROM scratch AS final
COPY --from=build /build/echo /echo
ENTRYPOINT ['/echo']