# syntax=docker/dockerfile:1
FROM golang:1.19

# Set destination for COPY
WORKDIR /app

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /gpa-calculator

EXPOSE 8080

# Run
CMD ["/gpa-calculator"]