# REST API in Go example
* gorilla/mux for routing
* lib/pq PostgreSQL driver
* docker/lib/postgres as database

___Source: https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql___

## Project structure
`app.go main.go model.go`

* `app.go` - App type to provide REST web service endpoints and model calls
* `main.go` - main package to run service
* `model.go` - defines products database model and executes CRUD operations against database

## CRUD of products
* create a new product: POST to /product
* delete an existing product: DELETE to /product/{id}
* update an existing product: PUT to /product/{id}
* fetch an existing product: GET to /product/{id}
* fetch a list of all existing products: GET to /products

## Setup Development Environment

### Create project workspace

```
$ cd ~/dev/go/github.com/amundsenjunior/rest-go-mux-pq
$ git clone https://github.com/amundsenjunior/rest-go-mux-pq.git
```

### Get project dependencies

You can pull the two dependencies directly, via:

```
$ go get github.com/gorilla/mux github.com/lib/pq
```

Or, using `go dep` (`go get -u github.com/golang/dep/cmd/dep`), use the present `Gopkg.*` files:

```
$ dep ensure
```

### Start PostgreSQL instance and create database & table
* https://hub.docker.com/_/postgres/
* https://www.postgresql.org/docs/9.2/static/app-psql.html
* https://www.tutorialspoint.com/postgresql/postgresql_create_database.htm

```
$ docker run --name rgmp -d -e POSTGRES_PASSWORD=restgomuxpq -p 5432:5432 postgres:alpine
$ docker run -it --rm --link rgmp:postgres postgres:alpine psql -h postgres -U postgres -W
  password: restgomuxpq
=# \list
=# CREATE DATABASE rgmp;
=# \c rgmp
  password: restgomuxpq
=# CREATE TABLE products (
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
  );
=# \dt
=# \q
```

### Testing

* TestEmptyTable - expected is a 200 code, but response returns 404
* TestGetNonExistentProduct - a GET request for a product by id should return error
* TestCreateProduct - an OK response code should come from a POST request of a new product
* TestGetProduct - an OK response code should come from a GET request on an existing product
* TestGetAllProducts - an OK response code and correct length of response body shoudl come from a GET request on all existing products
* TestUpdateProduct - an OK response code should come from a PUT request on an existing product to update it
* TestDeleteProduct - an OK response code should come from a DELETE request on removing an existing product


Execute TestMain with env vars:

```
$ export TEST_DB_USERNAME=postgres TEST_DB_PASSWORD=restgomuxpq TEST_DB_NAME=docker:5432/rgmp TEST_DB_SSLMODE=disable; go test -v
```

Alternatively, run `test_exec.sh` to start a Docker test database, execute the tests, and cleanup:

```
$ bash ./test_exec.sh
```

### Build and run the application

```
$ go build
$ export APP_DB_USERNAME=postgres APP_DB_PASSWORD=restgomuxpq APP_DB_NAME=docker:5432/rgmp APP_DB_SSLMODE=disable; ./rest-go-mux-pq
$ curl - POST -H "Content-Type: application/json" -d '{"name": "toy gorilla", "price": "29.99"}' http://localhost:8080/product
$ curl -X GET http://localhost:8080/products
```

## TODO
* doesn't error out when connection fails
* doesn't log address on which it's available

## GoDoc
* database/sql: https://golang.org/pkg/database/sql/
* gorilla/mux: http://www.gorillatoolkit.org/pkg/mux
* testing: https://golang.org/pkg/testing/
* net/http: https://godoc.org/net/http

## References
* database/sql tutorial: http://go-database-sql.org/
* JSON encoding/decoding: https://blog.golang.org/json-and-go
* when running a db unencrypted, include 'sslmode=disable' in connection string: https://stackoverflow.com/questions/21959148/ssl-is-not-enabled-on-the-server
* don't do local/relative imports: https://stackoverflow.com/questions/30885098/go-local-import-in-non-local-package
* when using 'go run' of multiple-file package, name all files in command: https://stackoverflow.com/questions/28153203/golang-undefined-function-declared-in-another-file
