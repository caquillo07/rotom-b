# Start from golang base image
FROM golang:alpine as builder

# Add Maintainer info
LABEL maintainer="Hector & Milla"

# Make sure to run `go mod vendor` before building the docker

# Copy the source from the current directory to the working Directory inside
# the container
WORKDIR /build
COPY . .

# Build the Go app
RUN GOOS=linux go build -o den-bot

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/data ./data
COPY --from=builder /build/den-bot .
COPY --from=builder /build/example-config.yaml config.yaml

RUN chmod +x den-bot

#Command to run the executable
CMD [ "./den-bot", "bot" ]
