# # Start from a Debian image with the latest version of Go installed
# # and a workspace (GOPATH) configured at /go.
# FROM golang

# # Copy the local package files to the container's workspace.
# ADD . /go/src/postGres

# WORKDIR /go/src/postGres

# # Build the outyet command inside the container.
# # (You may fetch or manage dependencies here,
# # either manually or with a tool like "godep".)
# # RUN go get github.com/lib/pq
# RUN CGO_ENABLED=0 GOOS=linux go install

# # Run the outyet command by default when the container starts.
# ENTRYPOINT /go/bin/postGres

# # Document that the service listens on port 8080.
# EXPOSE 8000

# # FROM golang:latest

# # WORKDIR /src
# # RUN mkdir /api
# # COPY ./ ./
# # RUN CGO_ENABLED=0 GOOS=linux go build -o /api/server
# # EXPOSE 8080
# # CMD ["/api/server"]

# Base build image
# FROM golang:1.11-alpine as base
FROM golang:latest

# Install some dependencies needed to build the project
# RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR /go/src/goTodos

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

# This is the ‘magic’ step that will download all the dependencies that are specified in 
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download 
# command will _ only_ be re-run when the go.mod or go.sum file change 
# (or when we add another docker instruction this line)
RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon


# This image builds the weavaite server
# FROM base as builder
# Here we copy the rest of the source code
COPY . .
# And compile the project
# RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -a -tags netgo -ldflags '-w -extldflags "-static"' .
# RUN CGO_ENABLED=0 go build -o /app .


EXPOSE 8000

#In this last stage, we start from a fresh Alpine image, to reduce the image size and not ship the Go compiler in our production artifacts.
# FROM alpine as app
# We add the certificates to be able to verify remote weaviate instances
# RUN apk add ca-certificates
# Finally we copy the statically compiled Go binary.
# COPY --from=builder /go/bin/goTodos /bin/goTodos
# ENTRYPOINT ["/app"]
ENTRYPOINT CompileDaemon -log-prefix=false -pattern="(.+\.go|.+\.c|.+\.html)$" -build="go build -o app ." -command="./app"

