FROM golang:alpine

# Set the Current Working Directory inside the container
WORKDIR /go/src/app/

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go mod download -x

# Install compile daemon for hot reloading
RUN go install -mod=mod github.com/githubnemo/CompileDaemon

# Expose port 8082 to the outside world
EXPOSE 8082

# Command to run the executable
ENTRYPOINT CompileDaemon --build="go build main.go" --command="./main"