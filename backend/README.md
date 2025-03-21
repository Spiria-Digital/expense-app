# Expense Manager Backend

This is the backend for the Expense Manager application. It is built using Golang, Gin & Bun framework.

## Running the app using docker
This is the **recommended way of running the app** if you simply want to integrate the frontend.
There is a docker file available in the backend directory, be sure to navigate to the backend directory before running the following commands.

You can run the following command to build the docker image:
```bash
docker build -t expense-manager-backend .
```

Once the image is built, you can run the following command to start the container:
```bash
docker run -p 8080:8080 expense-manager-backend
```
The backend will be running on `http://localhost:8080`.
You should be able to access the swagger documentation at `http://localhost:8080/api/swagger/index.html`.

## Setup for local development
Install go 1.23 or later from [here](https://golang.org/dl/). 
Navigate to the backend directory and run the following command to install the dependencies:
```bash
go mod download
```

### Database setup
The backend uses a SQLite database. You will need to compile the bun binary to run the database migrations. 
Run the following command to compile the bun binary:
```bash
go build -o bun `cmd/bun/main.go`
```

Once that's compiled, you will be able to run the database migrations:
```bash
./bun db init
./bun db migrate
```

After the migrations are complete, you can run the server:
```bash
go run main.go
```

The backend will be running on `http://localhost:8080`. You should be able to access the swagger documentation at `http://localhost:8080/api/swagger/index.html`.

## Generating Swagger Documentation
To generate the swagger documentation, you must install [swag](https://github.com/swaggo/swag) first. 
Run the following command from backend directory to generate the documentation:
```bash
swag init --parseInternal --parseDependency --parseDepth 2 --output server\docs
```