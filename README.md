# Xyfin 


## Prerequisites

- Docker
- PostgreSQL
- Go
- `migrate` tool for database migrations
- `sqlc` for Go code generation from SQL queries
- Protocol Buffers (`protoc`) for gRPC

## Setup

### 1. Network Setup

First, create a Docker network for the application:

```bash
make network
```

### 2. Database Setup

To set up the PostgreSQL database, run:

```bash
make postgres
```

After the database container is running, create the `xyfin` database:

```bash
make createdb
```

You can drop the database using:

```bash
make dropdb
```

### 3. Running the Server

To start the Xyfin server, use:

```bash
make server
```

### 4. Database Migrations

To apply all database migrations, run:

```bash
make migrateup
```

To apply a single migration, use:

```bash
make migrateup1
```

To rollback all migrations:

```bash
make migratedown
```

To rollback a single migration:

```bash
make migratedown1
```

### 5. Code Generation

Generate Go code from SQL queries using:

```bash
make sqlc
```

### 6. Protocol Buffers and gRPC

To regenerate Protocol Buffers and gRPC code:

```bash
make proto
```

### 7. Testing

Run the test suite with:

```bash
make test
```

### 8. Mock Generation

Generate mock implementations for testing:

```bash
make mock
```

### 9. Additional Services

Start a Redis container:

```bash
make redis
```

## Start Command

To start the Xyfin server with the necessary database connection, use the following command:

```bash
make run
```

This command starts the previously created Docker container for PostgreSQL and runs the Xyfin application.

## Documentation

To build database documentation:

```bash
make db_docs
```

To generate a SQL schema from the DBML:

```bash
make db_schema
```

## Development Tools

- **Evans CLI**: gRPC CLI client for testing your gRPC services.

```bash
make evans
```

## Contributing

1. Fork the repository.
2. Create a feature branch.
3. Commit your changes.
4. Push to your branch.
5. Create a new Pull Request.

## License

Xyfin is licensed under the MIT License. See [LICENSE](LICENSE) for more information.