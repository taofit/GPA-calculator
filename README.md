# GPA-calculator

## Run program

- Run migration up command: `make migrate-up` to perform the database migration with database schema
- Start docker containers command: `make start`
  - it will run `docker compose up`
  - run api server
  - seed the database table
  - run the tests
- Query the api endpoint, go to postman or insomnia enter url `localhost:8080/students/gpa` and see the result

## Code explain

- Migration folder is for database migration and it is needed by https://github.com/golang-migrate/migrate
- Api folder contains code that acts like handler or controller for the api call
- Database folder is where the sql queries are

## Migration Command

### Add more migration

If more tables need to be added.

- Create migration command for table 1: `make migrate-create name=create_table_1_name` # create_grade_table
- Create migration command for table 2: `make migrate-create name=create_table_2_name` # create_grade_scale_table
- Or put all create tables sql in one migration

### Migrate to previous state

- Running migration down command: `make migrate-down`

### Alternatively run migration without using make file

- Create migration command: `docker compose -f docker-compose.yml --profile tools run --rm migrate create -ext sql -dir /migrations create_grade_scale_table`
- Running migration up command: `docker compose -f docker-compose.yml --profile tools run --rm migrate up`
- Running migration down command: `docker compose -f docker-compose.yml --profile tools run --rm migrate down`

## Further work

- Add cache or pagination in the GPA end point, as it will fetch all students GPAs
- Add students grade and score to via an end point
- Add more test function and test cases
- Or may seed the in another server

## Thanks
