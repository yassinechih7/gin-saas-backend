FROM golang:1.19.4-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk add --no-cache bash git openssh g++ gcc clang

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# cp dot env
RUN mv .env_rename_me .env

# Build the Go app
RUN go build -o gin-boilerplate .

# Expose port 9000 to the outside world
EXPOSE 9000


# Run the executable
CMD ["./gin-boilerplate"]
