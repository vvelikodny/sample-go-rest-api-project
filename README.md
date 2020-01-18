[![Actions Status](https://github.com/vvelikodny/sample-go-rest-api-project/workflows/build/badge.svg?branch=master)](https://github.com/vvelikodny/sample-go-rest-api-project/actions)

## Getting Started

```shell
# checkout the project
git clone git@github.com:vvelikodny/sample-go-rest-api-project.git

cd sample-go-rest-api-project

# fetch deps
go mod vendor
```

## Run through docker-compose (recommend)
```
# docker-compose up -d
```

## Run through make
```
# start a PostgreSQL database server in a Docker container
make db-start

# migrate DB
make migrate

# run the RESTful API server
make run
```

At this time, you have a RESTful API server running at `http://127.0.0.1:3000/v1/`.

## REST API Tests, DB should be started

```
# run tests
make test
```
