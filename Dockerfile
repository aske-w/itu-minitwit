# Use go image with gcc compiler
FROM golang:1.17-stretch

# Create working dir
RUN mkdir /app

# Copy all files locally from the Docker file to our working direcotry
ADD ./server /app

# Set working dir
WORKDIR /app

# Build the project and run it
RUN go build -o main .
CMD ["/app/main"]