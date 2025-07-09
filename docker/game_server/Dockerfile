# Stage 1: Build the Go application
FROM golang:1.24 as builder

WORKDIR /game_server

COPY ../../go.mod go.sum ./

RUN go mod download

COPY ../.. .

RUN ./bin/build_game_server.sh

# Stage 2: Create a smaller image with the compiled binary
FROM debian:stable

WORKDIR /game_server

COPY --from=builder /game_server/target/game_server .

RUN chmod 777 server

EXPOSE 8000

CMD ["./game_server"]