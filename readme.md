# Idiomatic Go Todo Store

This is a more idiomatic-go approach to building a basic todo store, compared to my [other repo](https://github.com/grantjames/golang-todo-app) which follows a more object oriented apporach.

## Building and running

Build and run the CLI by doing the following

```
cd cmd/cli
go build
./cli -h
```

Passing the -h flag to the CLI will tell you what flags can be passed to it to interact wiht it.

Build and run the API by doing the following

```
cd cmd/api
go build
./api
```

## Design

The CLI and API are both small, single file applications that both call to the store in the todos.go file.

This store uses the actor pattern to ensure safe concurrent read and write to the map of todos that it stores.

## Tests

There is a parallel test in `todos_test.go` that runs concurrent creates and then deletions and then asserts that the number of todos in the store matches what is expected.