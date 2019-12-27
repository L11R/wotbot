# build the binary
FROM golang:1.13 AS build

COPY go.mod go.sum /
RUN go mod download

RUN useradd -u 10001 gopher

WORKDIR /go/src/wotbot

COPY . .

RUN GOOS=linux GOARCH=amd64 make build

# run the binary
FROM alpine:latest

COPY --from=build /etc/passwd /etc/passwd

USER gopher

COPY --from=build /go/src/wotbot/migrations /migrations
COPY --from=build /go/src/wotbot/bin/wotbot /wotbot

CMD ["/wotbot"]