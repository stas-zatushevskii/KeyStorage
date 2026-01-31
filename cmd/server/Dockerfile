# --- build ---
FROM golang:1.25 AS builder
WORKDIR /src/cmd/server

# копируем файлы модуля
COPY cmd/server/go.mod ./
# go.sum у тебя, судя по скрину, лежит в корне репо (это плохо для модуля).
# Можно пока не копировать go.sum — он создастся сам при download/tidy.
RUN go mod download

# копируем код модуля
COPY cmd/server/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server .

# --- run ---
FROM alpine:3.19
WORKDIR /server
RUN apk add --no-cache ca-certificates

COPY --from=builder /out/server /server/server

EXPOSE 8080
ENTRYPOINT ["/server/server"]