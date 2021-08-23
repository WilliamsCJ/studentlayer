FROM golang:1.16-alpine

# Set working directory
WORKDIR /app

# Copy go modules files and download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy go program files and build binary
COPY *.go ./
COPY data/universities_clean.json ./data/
RUN go build -o /studentlayer
EXPOSE 8080

# Run binary
CMD [ "/studentlayer" ]
