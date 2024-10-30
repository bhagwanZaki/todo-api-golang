# Basic Todo Api 

## Requirements
- Go (>=1.22)
- Postgres 17

## Installation

To install dependencies, run the following command:
```cmd
    go mod tidy
```

## Run server

### Using air (recommended)
if you have <a href="https://github.com/air-verse/air">air</a> package download
```cmd
    air
```
### Without air
Alternatively, you can run the server directly using:
```cmd
    go run main.go
```

## Setup

All SQL queries to set up tables, functions, and procedures can be found in the sql folder. Follow these to configure the database properly.