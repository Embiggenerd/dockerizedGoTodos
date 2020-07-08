

FROM golang:latest

WORKDIR /go/src/goTodos

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

# This step that will download all the dependencies that are specified in 
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download 
# command will _ only_ be re-run when the go.mod or go.sum file change 
# (or when we add another docker instruction this line)
RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon

COPY . .

EXPOSE 8000

# CompileDaemon gives us auto code reload for development
ENTRYPOINT CompileDaemon -log-prefix=false -pattern="(.+\.go|.+\.c|.+\.html)$" -build="go build -o app ." -command="./app"

