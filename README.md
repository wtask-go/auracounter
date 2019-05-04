# auracounter

* RPC HTTP server (`aurasrv`) to maintain a distributed cyclic counter with REST API.
Server binds counter ID, so you can run multiple instances to support single or several counters.

## API

Check API documentation and examples at https://documenter.getpostman.com/view/6496185/S1EJWgGQ

## Prerequisites

1. Install Go for your platform.

2. It is highly desirable to install a simple and convenient `godotenv` utility to run applications with a given set of environment variables:

```
> go install github.com/joho/godotenv
```

3. Clone or download application repository https://github.com/wtask-go/auracounter


> You may use any local directory as far as `auracounter` is go-module

4. Install docker and docker-compose

## Running tests

Move into project root and start project dependencies:

```
> godotenv -f .\deployments\config.test.env docker-compose -f .\deployments\docker-compose.yml up -d
```

After all dependencies will be started, run tests (all, including integrations):

```
> godotenv -f .\deployments\config.test.env go test  ./... -tags integration -v
```

To remove testing environment, run:

```
> godotenv -f .\deployments\config.test.env docker-compose -f .\deployments\docker-compose.yml down
```

## Running in dev-environment

Start dev-environment from project root:

```
> godotenv -f ./deployments/config.dev.env docker-compose -f ./deployments/docker-compose.yml up -d
```

> At the first time, you should wait until dependencies  will start. You may track a progress by checking docker-compose logs.

Run `aurasrv` in console (press Ctrl+C to stop server):

```
> godotenv -f ./deployments/config.dev.env go run ./cmd/aurasrv/.
```

You should see something like this:

```
aurasrv [2019-04-10 22:28:08.469356] INFO Server is starting ...
aurasrv [2019-04-10 22:28:08.469356] INFO Server is ready!
aurasrv [2019-04-10 22:28:08.469356] INFO Server has stopped, bye ( ᴗ_ ᴗ)
```

When the server is running you can use [Postman](https://www.getpostman.com/) or other http-client to work with server API. Check [API docs](https://documenter.getpostman.com/view/6496185/S1EJWgGQ).

Also, you can get help from the server in console:

```
> go run ./cmd/aurasrv/. -h
```

### Stopping dev-environment

If you want to stop/start server environment fast, run:

```
> docker-compose -f .\deployments\docker-compose.yml stop
```

Start dev-environment again.

```
> docker-compose -f .\deployments\docker-compose.yml start
```

Completely remove all used containers from your host and stop dev-environment:

```
> docker-compose -f .\deployments\docker-compose.yml down
```

# Feature plans

* Add support for OPTIONS method (HTTP) to expose API
* Add support for `make`
* Add test for whole server

...
