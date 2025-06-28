# Stage 1: Build the Go application
FROM golang:1.24 as builder

WORKDIR /lastmudserver

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN ./bin/build.sh

# Stage 2: Create a smaller image with the compiled binary
FROM debian:stable

WORKDIR /lastmudserver

COPY --from=builder /lastmudserver/target/lastmudserver .

RUN chmod 777 lastmudserver

EXPOSE 8000

CMD ["./lastmudserver"]