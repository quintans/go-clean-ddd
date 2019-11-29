# Example of Clean Architecture + DDD in go

## Architecture
http://dddsample.sourceforge.net/architecture.html

## Characterization
http://dddsample.sourceforge.net/characterization.html

## Clean Code
https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

Start the database container running the following command in the project root:

```sh
docker run --name postgres -e POSTGRES_PASSWORD=secret -p 5432:5432 -v "$PWD/scripts":/docker-entrypoint-initdb.d -d postgres:9.6.8
```
