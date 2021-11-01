# A CRUDy webapp for Albums or something

A re-imagining of the Go tutorial: https://golang.org/doc/tutorial/web-service-gin

## Running

It's golang, so it's really easy:

    go get .

    go run main.go

The client will run in your browser at http://localhost:8000

All data is loaded using ajax calls because it's more fun. 

## Using API

The RESTy API is pretty standard and works exactly like you think it will: 

### Get all records

    curl http://localhost:8000/albums

### Get record by ID

    curl http://localhost:8000/album/1

### Delete record by ID

    curl -X DELETE http://localhost:8000/album/1 

### Create record
    curl -X POST http://localhost:8000/album/ -d '
        {
            "title":"Shrek 2: The Soundtrack",
            "artist":"Dreamworks",
            "price":5.99
        }'
