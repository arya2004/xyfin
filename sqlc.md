# SQLC User Guide

SQLC is a tool that generates Go code from SQL queries. This guide provides an overview of all SQLC commands and their usage.

## Installation

To install SQLC, run:
```sh
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
```

## Commands

### 1. `sqlc init`
Initializes a new SQLC project by creating a `sqlc.yaml` configuration file in the current directory.

#### Usage:
```sh
sqlc init
```

### 2. `sqlc generate`
Generates Go code from the SQL queries and schema defined in the `sqlc.yaml` configuration file.

#### Usage:
```sh
sqlc generate
```

### 3. `sqlc compile`
Compiles the SQL queries and prints any parsing errors without generating Go code. This is useful for validating your SQL queries.

#### Usage:
```sh
sqlc compile
```

### 4. `sqlc vet`
Lints the SQL queries to check for common mistakes and best practices. This command helps ensure your SQL queries are optimized and follow good conventions.

#### Usage:
```sh
sqlc vet
```

### 5. `sqlc completion`
Generates shell completion scripts for bash, zsh, fish, and PowerShell. This command helps with auto-completion of SQLC commands in the terminal.

#### Usage:
```sh
sqlc completion [bash|zsh|fish|powershell]
```

### 6. `sqlc version`
Prints the current version of SQLC.

#### Usage:
```sh
sqlc version
```

## Configuration

SQLC uses a `sqlc.yaml` file to configure the code generation process. Below is an example configuration file:

```yaml
version: "1"
packages:
  - path: "db"
    queries: "./sql/queries.sql"
    schema: "./sql/schema.sql"
    engine: "postgresql"
```

### Configuration Options

- `version`: The configuration file version. Always set this to `"1"`.
- `packages`: A list of packages to generate.
  - `path`: The output directory for the generated Go code.
  - `queries`: The path to the file containing SQL queries.
  - `schema`: The path to the file containing the database schema.
  - `engine`: The SQL database engine (e.g., `postgresql`, `mysql`, `sqlite`).

## Example Workflow

1. **Initialize a new SQLC project**:
    ```sh
    sqlc init
    ```

2. **Define SQL queries** in `queries.sql`:
    ```sql
    -- name: GetUserByID :one
    SELECT id, name, email FROM users WHERE id = $1;
    ```

3. **Define the database schema** in `schema.sql`:
    ```sql
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL
    );
    ```

4. **Configure SQLC** in `sqlc.yaml`:
    ```yaml
    version: "1"
    packages:
      - path: "db"
        queries: "./sql/queries.sql"
        schema: "./sql/schema.sql"
        engine: "postgresql"
    ```

5. **Generate Go code**:
    ```sh
    sqlc generate
    ```

6. **Use the generated code** in your Go application:
    ```go
    package main

    import (
        "context"
        "log"
        "your_project/db"
    )

    func main() {
        // Assuming db.Conn is your database connection
        q := db.New(db.Conn)
        user, err := q.GetUserByID(context.Background(), 1)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("User: %+v\n", user)
    }
    ```

