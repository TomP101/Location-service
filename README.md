# Location-service

## Project Overview

This project is a microservices-based system in Go for managing and tracking user locations.  
It consists of two services:

1. **Location Service** – Handles location updates and user search by proximity.  
2. **Location History Service** – Records historical data and calculates the distance traveled by users over a period of time.

## Features

- **Update User Location** (`POST /locations`)  
- **Search Users by Location** (`GET /search`)  
- **Calculate Distance Traveled** (`GET /history/distance`)  


## Technologies Used

- **Go (Golang)**  
- **Gin** for REST API  
- **gRPC** for service communication  
- **MySQL** for data storage

## Local Setup Instructions

### 1. Clone the Repository

git clone https://github.com/TomP101/location-service.git  
cd location-tracking-system  

### 2. Set Up Environment Variables:

export DBUSER=root
export DBPASS=

### 3. Install Dependencies
Download and install Go from the official site: https://golang.org/dl/
Install mysql:
  brew install mysql  # For macOS  
  sudo apt-get install mysql-server  # For Linux
  Start MySQL:
    brew services start mysql  # For macOS  
    sudo service mysql start  # For Linux  
Install GRPC and Protobuf
  brew install protobuf  # For macOS  
  sudo apt-get install protobuf-compiler  # For Linux 
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest  
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

run go mod tidy

### 4. Set Up MySQL Database

After logging into your mysql terminal run :
source path_to_your_repository/create-tables.sql

this will create the required tables for the database

### 5. Run Services(using terminal)

Location service:
  cd location-service  
  go run main.go  
Location History service:
  cd location-history-service  
  go run main.go  

### 6. Test the app
After starting both services run these commands in another terminal to check the application:
1. Add a new user:
curl -X POST http://localhost:8080/locations \
-H "Content-Type: application/json" \
-d '{"name":"test_user","latitude":35.12314,"longitude":27.64532}'
2. Update location for the same user:
curl -X POST http://localhost:8080/locations \
-H "Content-Type: application/json" \
-d '{"name":"test_user","latitude":39.12355,"longitude":27.64538}'
3. Retrieve all locations:
curl http://localhost:8080/locations
4. Calculate distance traveled by user
curl "http://localhost:8081/history/distance?username=test_user"
5. Search for Users within a Radius:
curl "http://localhost:8080/search?latitude=35.0&longitude=27.0&radius=5000"


