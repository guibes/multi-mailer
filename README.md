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

## TODO

- [ ] Sendgrid integration
- [ ] List another email clients
- [ ] GUI
- [ ] CLI

