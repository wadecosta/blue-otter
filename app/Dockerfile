# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory in the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
RUN go build -o main .

# Install MySQL
RUN apt-get update && apt-get install -y mysql-server

# Expose the MySQL port
EXPOSE 3306

# Start the MySQL service
CMD ["service", "mysql", "start"]

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

# Create a MySQL database named "yourdbname"
RUN service mysql start && \
    mysql -u root -e "CREATE DATABASE yourdbname;"

