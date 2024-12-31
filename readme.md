# Todo Api

## Overview
A basic Todo API built with Go, featuring secure authentication and RabbitMQ integration for message queuing.

---

## Requirements

- Go (>=1.22)
- Postgres 17
- RabbitMQ server

---

## Installation

To install dependencies, run the following command:
```bash
    go mod tidy
```

---

## Run Server

```bash
    cd todo
```
And then use one of the following to run the server
### Using Air (Recommended)
If you have the <a href="https://github.com/air-verse/air">Air</a> package installed:
```bash
    air
```

### Without Air
Alternatively, you can run the server directly using:
```bash
    go run main.go
```
# Run the consumer

```bash
    cd consumer
    go run .
```
---

## Setup

- All SQL queries to set up tables, functions, and procedures can be found in the `sql` folder.
- Follow these instructions to configure the database properly.

---

## Features

- **Authentication**: Secure user login and registration.
- **RabbitMQ Integration**: Efficient handling of asynchronous operations.
- **PostgreSQL**: Robust relational database for storing todos and user data.

---

## Development Notes

- **RabbitMQ:**
  - Utilized for task/event queuing and asynchronous processing.

- **Database Setup:**
  - Ensure the PostgreSQL instance is correctly configured using the SQL scripts provided.

---

## Acknowledgments

- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
- [Go Programming Language](https://golang.org/doc/)
