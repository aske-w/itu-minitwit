# Use go image with gcc compiler
FROM golang:1.17-stretch


# Create working dir
RUN mkdir /app

# Copy
ADD ./server /app

# Set working dir
WORKDIR /app

RUN go build -o main .
CMD /app/main