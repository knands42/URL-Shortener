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

## Technical Decisions

This application was built thinking like any other application on an early stage would be built, but the real world solution for a highly available system is ilustrated in the diagram in the root of the project: `architecture.drawio`.

### Storage

The choice between SQL and NoSQL databases was made based on the fact that the application needs to be highly available, consistent and query efficient, so the best choice was to use a SQL database even though NOSQL can scale horizontally and is more flexible.

In the future more metadata will be added like user information and how they are accessing and generating new short urls, so a SQL database will help with the relationships between the tables, but this the current implementation no need to extract the most of the **3rd normal form**.

**RDS Aurora** would be one of the best choice for a production environment since it was build to scale horizontally.

### Caching

The application uses **Redis** to cache the short urls and the original urls, this is a good choice to avoid hitting the database every time a request is made.

A good strategy is to rely on the **LRU eviction policy** to avoid memory leaks and to keep the cache size under control if the cache is in-memory (which for a first time with a few users it may be good enough).

### Graceful Shutdown

In order to avoid memory leasks and to make sure that the application is not processing any request when it is shutting down or any connection remains open, the application uses the `os.Signal` to catch the `SIGTERM` signal and to gracefully shutdown the application, like closing the database, redis and jaeger connection.

### Observability

The application uses **Jaeger** to trace the requests and to have a better understanding of the application performance and to debug any issues that may arise (like `where` it did happen, not only `when`).

Besides that, the application uses **logs** (which is a must) middleware and custom logs for some failures that may happen in the application.

A centralized collector and the power of OTEL make the appplication work with many vendors and delegate the responsibility of managing what to do with the data that is being collected, like metrics, traces and logs.

> **Note**: Other observability tools that will benefit the application are: Prometheus for metrics, APM to manage mostly the health of the app and Grafana for a centralized place to connect many datasources and have nice dashboards visualization for the project.

### Security

For testing purposes the **cors** is enabled for all origins, but in a production environment it should be restricted to the domains that are allowed to access the application.

### Background Jobs

Since the application is highly available, every time some use tries to access informations regading the short url or just delete it, the application needs to record this actions somehow.

Instead of doing synchronously and hang the TCP connection for each GET request just to process this actions, for a production environment the so called **outbox pattern** would be a good choice to avoid this kind of problem by relying on the existing database to store this actions (events) and then a background job would process this actions asynchronously (this will work because this is a long running process and not serverless) no **broker** is needed (unless we start to implement connections through microservices).

> **Note**: Another reason to use a SQL database is to exploit the **ACID** features that it provides (talking about consistency here), so background jobs can rely on it for this outbox pattern

### Resilience

The application does not contain any resilience pattern manually implemented, but a **circuit breaker** or **rate limiter (per user)** would help to avoid cascading failures.
