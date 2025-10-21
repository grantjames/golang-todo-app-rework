# Idiomatic Go Todo Store

This is a more idiomatic-go approach to building a basic todo store, compared to my [other repo](https://github.com/grantjames/golang-todo-app) which follows a more object oriented apporach.

The solution meets the requirements detailed in the `requirements.md` file.

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

### Todo store

There are two implementations of the todo store. The one used by the CLI and API encapsulates the functionality inside a file and just uses unexported variables to ensure ownership of data.

The other implementation is in `todos_struct.go` (commented out). This works similarly to the "no struct" version, but has a struct own its own data and by instance methods on it for consumers to interact with it.

I'm unsure of the "best" approach here, but thought I'd try out both.

The biggest change, however, from my previous solution, is that the store itself implements the actor pattern, rather than having an Actor wrapping round a store. This way, the actor can just export methods like `Create()` and `Get()` which send messages to it's own "queue"/channel, rather than in my previous solution where the actor just had one generic method to recieve "messages". This worked well but required lots of "transport" structs to hold a request and a response type for every type of interaction with the store.

## Tests

There is a parallel test in `todos_test.go` that runs concurrent creates and then deletions and then asserts that the number of todos in the store matches what is expected.