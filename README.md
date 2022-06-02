# go-clean-ddd 
Example of Clean Architecture + DDD + CQRS in go

> There are many ways to structure a project that follows the above architecture.
> 
> The most important thing is to keep the business logic separated from the technological details and to respect the dependency hierarchy between the different layers.

Start the database container by doing

```sh
docker run --name postgres -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:9.6.8
```

or

```
task up
```

> need [Task](https://taskfile.dev/installation/) installed

## Infra
this is where all the technological details reside.

- Database connection
- Web Server bootstrap

## Controller
In the [controller](./internal/controller/) folder is where we handle requests from the outside of our domain, interpret and transform it into an model known to the domain and then forward the call to the proper use case.

> eg: convert a json payload into a command and return an ID representing the side effect

## Gateway
The [gateway](./internal/gateway/) we handle calls to external services. We receive the domain object, transform it to the specific tech and return any reply by doing the reverse process.

> eg: calling a database with an aggregate ID and return an aggregate.

## Use Cases
This is were the application rules are. Here we will have the services needed to satisfy requirements from an external call. Can be from another microservice, from web or mobile.
Depending on the caller, usually each should have its own method because they would represent distinct use cases.

Most of the time, the implementation just retrieves an aggregate (where the business rules live) from the repository calls a method to apply the business rules and saves the aggregate.
There are more complex scenarios where a third party service may need to be called before or after applying the business rules.

Here we also apply the simplest form of **CQRS**, a pattern that separates read and update operations.

### Interfaces
Here we also find all the interfaces for all the use cases (incoming calls) and gateways (outgoing calls).
The inputs and outputs declared in the interfaces should never refer to a specific technology. You will not find in the structs definition any tag like ``` `json:"..."` ``` or ``` `\db:"..."``` `

When naming an interface we avoid referring to a specific technology. For example, KafkaPublisher is a bad name, since it refers to specif technology. We should indicate __intent__, for example, if what we want is to notify something, __notifier__ would be a better name for the interface. It is hiding the implementation details. The implementation can push to a kafka topic, write in a DB or send an email. The domain does not care.


### command
thing to do and it will produce a side effect. Theoretically it shouldn't return anything, since it is not a Query, but lets be pragmatic and return an output when needed.

It is here that we use the repository to get the aggregate, call the use case, and save the aggregate. If there is a need to propagate Domain Event, the outbox pattern could be used (if we don't have another way to propagate database changes). All is done inside the same database transaction. 

### query
get information in the form of DTOs

Both [command](./internal/usecase/command/) and [query](./internal/usecase/query/) use cases will have to implement the interface defined in [ports](./internal/usecase/)


## Domain
Where the enterprise business rules are. It uses aggregates to enforce all the business rules (invariants). By using an aggregate we can skip the use of Domain Services since we can always pass a Policy that wraps any external verification.
This is where the aggregates, entities, value objects and domain events live.


## Resources

### DDD

http://dddsample.sourceforge.net/architecture.html

http://dddsample.sourceforge.net/characterization.html

### Clean Code
https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

There is an interesting article that makes a [Comparison of Domain-Driven Design and Clean Architecture Concepts](https://khalilstemmler.com/articles/software-design-architecture/domain-driven-design-vs-clean-architecture/)
