#!/bin/bash

# Generate Swagger documentation for the API

echo "Installing swag if not present..."
go install github.com/swaggo/swag/cmd/swag@latest

echo "Generating Swagger documentation..."
swag init -g cmd/main.go -o docs --parseDependency --parseInternal

echo "Swagger documentation generated successfully!"
echo "Access the documentation at: http://localhost:3000/swagger/index.html" 