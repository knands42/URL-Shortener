# URL-Shortener

Just a simple go application to leverage my knowledge in Golang and System Design

## Execute the application

### Pre requisites

Before running the application, make sure you have the following golang tools installed with the following command:

> **Note**: Golang version 1.20 (or higher), Python3.10 (or higher) and docker is required to be already pre installed in your machine.

```sh
make setup
```

This will install swagger and migrate cli tools used by the application and blazemeter to run the performance tests.

### Boot up the application

Start the depencies with the following command:

```sh
docker compose up
```

Generate the swagger documentation:

```sh
make gen-docs
```

Migrate the database:

```sh
make migrate-up
```

Then run the application with the following command:

```sh
make build-and-run
```

### Useful links

- [Swagger Documentation](http://localhost:3333/swagger/index.html)
- [Healthcheck](http://localhost:3333/health)

### Run the tests

First check if integration tests are passing:

```sh
make integration-tests
```

Then execute the performance tests:

```sh
make performance-tests
```

## System Design
