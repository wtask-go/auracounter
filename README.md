# auracounter
The set of applications to maintain of distributed cyclic counter:

* REST-like RPC server to manage and deal with distributed counter. Server binds counter ID (cid), so you can run multiple instances to support single counter.
* CLI support may appear

# API
Check API documentation and examples at https://documenter.getpostman.com/view/6496185/S1EJWgGQ

# Prerequisites

1. Used GO version:

```
go version go1.12.1 windows/amd64
```

2. It is highly desirable to install a simple and convenient `godotenv` utility to run applications with a given set of environment variables:

```
> go install github.com/joho/godotenv
```

3. Clone or downlod application repository https://github.com/wtask-go/auracounter


> You may use any local directory as far as `auracounter` is go-module

4. Install docker and docker-compose

# Runing in dev-environment

Use previously installed `godotenv` to start dependencies and to run REST-server.

1. Dive into project root and start server environment:

```
> godotenv -f ./deployments/config.dev.env docker-compose -f ./deployments/docker-compose.yml up -d
```

At first time, you should wait until MySQL container will start. You may check progress by open docker-compose logs.

2. Run RPC-server in console, press Ctrl+C to stop server:

```
> godotenv -f ./deployments/config.dev.env go run ./cmd/aurasrv/.
```
You should see something like this:

```
aurasrv [2019-04-10 22:28:08.469356] INFO Server is starting ...
aurasrv [2019-04-10 22:28:08.469356] INFO Server is ready!
aurasrv [2019-04-10 22:28:08.469356] INFO Server has stopped, bye ( ᴗ_ ᴗ)
```
When the server is running you can use [Postman](https://www.getpostman.com/) or other http-client to work with server API. Check [API docs](https://github.com/wtask-go/auracounter).

Also, you can get help from the server in console:

```
> go run ./cmd/aurasrv/. -h
```

3. Stop environment

If you want to stop/start server environment fast, run:

```
> docker-compose -f .\deployments\docker-compose.yml stop
```

or

```
> docker-compose -f .\deployments\docker-compose.yml start
```

Or to remove all used containers from your host:

```
> docker-compose -f .\deployments\docker-compose.yml down
```

# Feature plans

* Log server requests and errors with logging.Facage
* Add support for OPTIONS method (HTTP) to expose API
* Add support for `make`
* Add integration test for MySQL-repository
* Add test for whole server
* Add missing docs and fix outdated comments or typos

...
