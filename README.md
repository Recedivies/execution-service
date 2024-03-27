<div align="center" style="padding-bottom: 20px">
  <h1>Execution Service</h1>
</div>

## Setup local development

### Install tools

-   [TablePlus](https://tableplus.com/)
-   [Golang](https://golang.org/)
-   [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

-   [DB Docs](https://dbdocs.io/docs)

    ```bash
    npm install -g dbdocs
    dbdocs login
    ```

-   [DBML CLI](https://www.dbml.org/cli/#installation)

    ```bash
    npm install -g @dbml/cli
    dbml2sql --version
    ```

-   [Sqlc](https://github.com/kyleconroy/sqlc#installation)

-   [Gomock](https://github.com/golang/mock)

    ```bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

## Setup infrastructure

### Documentation

-   Generate DB documentation:

    ```bash
    make db_docs
    ```

### How to generate code

-   Generate schema SQL file with DBML:

    ```bash
    make db_schema
    ```

-   Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```

-   Create new db migration:

    ```bash
    migrate create -ext sql -dir db/migration -seq <migration_name>
    ```

### How to run

-   Run server:

    ```bash
    make server
    ```
