# GPA-calculator

## Run program

- Running migration up command: `make migrate-up` to perform the database migration with database schema
- Start docker containers command: `make start`

## Migration Command

### Add more migration

If more tables need to be added.

- Create migration command for table 1: `make migrate-create name=create_table_1_name` # create_grade_table
- Create migration command for table 2: `make migrate-create name=create_table_2_name` # create_grade_scale_table
- or put all create table sql in one migration

### Migrate to previous state

- Running migration down command: `make migrate-down`

### Alternatively run migration without using make file

- Create migration command: `docker compose -f docker-compose.yml --profile tools run --rm migrate create -ext sql -dir /migrations create_grade_scale_table`
- Running migration up command: `docker compose -f docker-compose.yml --profile tools run --rm migrate up`
- Running migration down command: `docker compose -f docker-compose.yml --profile tools run --rm migrate down`
