# GPA-calculator

## Run program

- start docker containers command: `make start`
- Create migration command for table 1: `make migrate-create name=create_grade_table`
- Create migration command for table 2: `make migrate-create name=create_grade-scale_table`
- Running migration up command: `make migrate-up`
- Running migration down command: `make migrate-down`

## Alternatively run migration without using make file

- Create migration command: `docker compose -f docker-compose.yml --profile tools run --rm migrate create -ext sql -dir /migrations create_grade-scale_table``
- Running migration up command: `docker compose -f docker-compose.yml --profile tools run --rm migrate up`
- Running migration down command: `docker compose -f docker-compose.yml --profile tools run --rm migrate down`
