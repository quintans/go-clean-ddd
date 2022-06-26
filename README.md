# go-clean-ddd 
Example of Clean Architecture + DDD + CQRS in go

To demonstrate the of several architecture patterns a simple demo is provided.

## Demo
This demo simulates a customer registration.

It intends to demonstrate the following:
- project structure 
- [outbox pattern](https://microservices.io/patterns/data/transactional-outbox.html) to update database and send events
- integration event, calling an external gateway: fake message queue
- domain event for dependency (transaction) between 2 aggregates (synchronous)
- unit of work across 2 aggregates
- optimistic locking

> transaction between two or more aggregates is done through the use of domain events

The business rules  are the following:
- When a new customer wants to register, it will provide an email. If the email is unique a new registration aggregate is created.
- Then the customer will receive an email (not really) with an url that needs to be called to create the customer account
- After the successful activation the customer can complete the rest of the information

__How__:
- create a registration record if the email is unique
- When saving a new registration, a new event is stored in the outbox table, in the same transaction.
- A background process polls the new event and publishes to a message queue.
An event handler will then process this event.
In a production environment this event handler would be responsible to send am email to the customer with a confirmation link but for demonstration purposes we are going to fake the sending of the email and the customer customer confirmation by directly calling the confirmation endpoint as soon the registration is done
- an endpoint with the registration ID is called, and after a successful validation in the __registration__ aggregate a domain event is fired, leading to the creation of the customer aggregate and deletion of the registration.
All will be saved in the same transaction.

> For demonstration purposes the queue is fake

> There are many ways to structure a project that follows the above architecture. Each team should find the one that best suits them as long 
> it keeps the business logic (domain) separated from the technological details (infrastructure), respects the dependency hierarchy between the different layers (controller -> domain <- gateway) and respects the control flow (controller -> domain -> gateway)
>
> An effort is made to keep the packages as shallow as possible

Start the database container by doing

```sh
docker run --name postgres -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:9.6.8
```

or

```
task up
```

> need [Task](https://taskfile.dev/installation/) installed


The project structure separates technology, application and domain in 3 high level folders respectively named `infra`, `app` and `domain`.

## Infra
this is where all the technological details reside.

- environment variables
- Database connection
- Web Server bootstrap
- configuration
- adapters (controllers and gateways)

### Controller
The [controller](./internal/infra/controller/) handles requests from the outside of our application, interpret and transform it into an model known to the domain and then forward the call to the proper use case.
This should as thin as possible.
No logic should be found here, only data transformation.

> It is fine to use value objects at this layer but don't use entities as an input.

> eg: convert a json payload into a command and return an ID representing the side effect

### Gateway
The [gateway](./internal/infra/gateway/) handles calls to external services. It receives the domain object, transform it to the specific tech and return any reply by doing the reverse process.

> eg: calling a database with an aggregate ID and return an aggregate.

## Domain
It contains is business logic and your domain modeling.
It's independent of specific technologies like databases or web APIs.

Domain is split into separate aggregates/entities packages to demonstrate the interaction between to domain.
Most of the cases one aggregate is enough for the domain and in these cases we can put everything  inside the package domain.

### Modules (Customer, Registration)
This is where the **enterprise** business rules are. It uses aggregates to enforce all the business rules (invariants).
By using an aggregate we can skip the use of Domain Services since we can always pass a Policy that wraps any external verification.
This is where the aggregates, entities, value objects and domain events related to the domain module live.

## App / Use Cases
This is where the **application** rules are. Here we will have the services needed to satisfy requirements from an external call. Can be from another microservice, from web or mobile.
Depending on the caller, usually each should have its own method because they would represent distinct use cases.

Most of the time, the implementation just retrieves an aggregate (where the enterprise business rules live) from the repository calls a method to apply the business rules and saves the aggregate.
There are more complex scenarios where a third party service may need to be called before or after applying the business rules.

Here we also apply the simplest form of **CQRS**, a pattern that separates read and update operations.

### ports
Here we also find all the ports (interfaces) for all the controllers (incoming calls) and gateways (outgoing calls).
The inputs and outputs declared in ports should never refer to a specific technology. You will not find in the structs definition any tag like ``` `json:"..."` ``` or ``` `\db:"..."``` `

When naming an interface we avoid referring to a specific technology. For example, KafkaPublisher is a bad name, since it refers to specif technology. We should indicate __intent__, for example, if what we want is to publish something, __Publisher__ would be a better name for the interface. It is hiding the implementation details. The implementation can push to a kafka topic, write in a DB or send an email. The domain does not care about the technological details.


### Command
thing to do and it will produce a side effect. Theoretically it shouldn't return anything, since it is not a Query, but lets be pragmatic and return an output when needed.

It is here that we use the repository to get the aggregate, call the use case, and save the aggregate. If there is a need to propagate Domain Event, the outbox pattern could be used (if we don't have another way to propagate database changes). All is done inside the same database transaction. 

### Query
get information in the form of DTOs

Both [command](./internal/domain/usecase/command/) and [query](./internal/domain/usecase/query/) use cases will have to implement the interface defined in [usecase](./internal/domain/usecase/)



## Resources

### DDD

http://dddsample.sourceforge.net/architecture.html

http://dddsample.sourceforge.net/characterization.html

### Clean Code
https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

There is an interesting article that makes a [Comparison of Domain-Driven Design and Clean Architecture Concepts](https://khalilstemmler.com/articles/software-design-architecture/domain-driven-design-vs-clean-architecture/)
