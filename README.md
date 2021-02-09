# wyrm
---
# Development Environment
## Database setup guide
1. Run database script
    ```sh
        ./scripts/db-up.sh
    ```
2. Update ```.env``` file with database params
    ```sh
        // Example
        DB_HOST=localhost
        DB_USER=admin
        DB_PASSWORD=admin
        DB_NAME=dev
        DB_PORT=5432
        DB_SCHEMA_PATH=cmd/db-init/schema.sql
    ```
3. Run database initialization script
    ```sh
        ./scripts/db-init.sh
    ```
---