FROM golang:bullseye

COPY . ./simply-ws
WORKDIR ./simply-ws

RUN go test ./...
RUN go build -o ./dist/app ./cmd/main.go

ENTRYPOINT ["./dist/app", "./resources/sws.sqlite3"]