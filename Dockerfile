
FROM golang:1.19-alpine

WORKDIR /go/src/app

# Copy the main.go file into the container
COPY main.go .

# Install Air using go install
RUN go install github.com/cosmtrek/air@latest

# Expose the port on which the application will run
EXPOSE 8080

# Command to run your application using Air
CMD ["air"]
