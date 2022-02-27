# Use go image with gcc compiler
FROM golang:1.17-stretch


# Create working dir
RUN mkdir /app

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait
# Copy all files locally from the Docker file to our working direcotry
ADD ./server /app

# Set working dir
WORKDIR /app

# Build the project and run it
RUN go build -o main .
CMD /wait && /app/main