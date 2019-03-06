# Simple Microservices example

- Includes authentication and authorization as a microservice

## TODO:

- Auth should be based on OpenID Connect.

## Layout/services

- `auth`: Authentication and authorization.
- `echo`: Really simple service that can echo what you send to it.
    - TODO: right now it doesn't do a simple echo, it includes the `name` you send to it in the Hi message on `/hi`.

## Demo

Start the servers with docker-compose:

```
source .env
docker-compose up --build
```

First you have to authenticate yourself. You can do this by calling the `auth` api:

```
```