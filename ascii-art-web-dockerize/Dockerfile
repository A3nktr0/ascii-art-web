# syntax=docker/dockerfile:1

FROM golang:1.20

# Metadata
LABEL "description"="Docker image of ascii web project who run go server and display GUI in web browser in order to translate words graphical ASCII art representation"
LABEL version="1.0"
LABEL author="aguillen | qdelooze"
LABEL support-contact="delooze.quentin@gmail.com"


# Set workspace directory
WORKDIR /app

# Copy all files and download dependencies
ADD . ./
RUN go mod download

# Create an executable file of the program
RUN go build -o /ascii-art-web-dockerize cmd/main.go

# Run
CMD ["/ascii-art-web-dockerize"]
