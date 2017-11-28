# REST API in Go example
* gorilla/mux for routing
* lib/pq PostgreSQL driver
* docker/lib/postgres as database

___Source: https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql___

## CRUD of products
* create a new product: POST to /product
* delete an existing product: DELETE to /product/{id}
* update an existing product: PUT to /product/{id}
* fetch an existing product: GET to /product/{id}
* fetch a list of all existing products: GET to /products

create project workspace
`$ cd ~/dev/go/github.com/amundsenjunior/rest-go-mux-pq`

get Go dependencies
`$ go get github.com/gorilla/mux github.com/lib/pq`

Start PostgreSQL instance and create database & table
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

create application structure
`$ touch app.go main.go main_test.go model.go`

define App type to hold application
define main package to run service
define products database model

start writing tests with test database
include pre-testing database setup and post-testing database cleanup

execute TestMain with env vars
`$ export TEST_DB_USERNAME=postgres TEST_DB_PASSWORD=restgomuxpq TEST_DB_NAME=docker:5432/rgmp TEST_DB_SSLMODE=disable; go test -v`

1. add and execute TestEmptyTable, where expected is a 200 code, but response returns 404
1. add and execute TestGetNonExistentProduct, where a GET request for a product by id should return error
1. add and execute TestCreateProduct, where an OK response code should come from a POST request of a new product
1. add and execute TestGetProduct, where an OK response code should come from a GET request on an existing product
1. add and execute TestGetAllProducts, where an OK response code and correct length of response body shoudl come from a GET request on all existing products
1. add and execute TestUpdateProduct, where an OK response code should come from a PUT request on an existing product to update it
1. add and execute TestDeleteProduct, where an OK response code should come from a DELETE request on removing an existing product

execute application with env vars
```
$ export APP_DB_USERNAME=postgres APP_DB_PASSWORD=restgomuxpq APP_DB_NAME=docker:5432/rgmp APP_DB_SSLMODE=disable; go run main.go app.go model.go
$ go build; export APP_DB_USERNAME=postgres APP_DB_PASSWORD=restgomuxpq APP_DB_NAME=docker:5432/rgmp APP_DB_SSLMODE=disable; ./rest-go-mux-pq
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
