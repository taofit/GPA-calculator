# syntax=docker/dockerfile:1
FROM golang:1.19

# Set destination for COPY
WORKDIR /app
RUN mkdir "/build"

COPY . .
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

EXPOSE 8080

ENTRYPOINT CompileDaemon -build="go build -o /build/app" -command="/build/app"