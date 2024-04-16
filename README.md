# Multi Mailer

A tool to send email using multiple providers, implemented in cloud native.

## Running as development

- Start a postgres or use the `docker-compose` file on the project.
- Install go modules using `go mod download`
- Use .env file to configure your envs and export it `export $(cat .env | xargs)`
- Run the migrations `go run ./cmd/migrate`
- Run the project `go run ./cmd/api`

## Create a new migration

- Install [goose](https://github.com/pressly/goose)
- Run `cd /migrations && goose <migration-name> sql && cd ..`

## Generate Swagger Documentation

- Install [swag](https://github.com/swaggo/swag)
- Run the command `swag init -g cmd/api/main.go -o .swagger -ot yaml`

Will output the swagger file in `.swagger/swagger.yaml`

The docs can be previwed on [Swagger Editor](https://editor.swagger.io)

## TODO

- [ ] Sendgrid integration
- [ ] List another email clients
- [ ] GUI
- [ ] CLI
