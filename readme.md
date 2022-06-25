# go gin mongo, LOL!

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Gintonica is a learning project under construction with:
- GO
- MongoDb
- GIN
- Jwt
- RabbitMQ
- Swagger

## setup mongo and rabbitMQ
```bash
docker-compose up -d
```

## install db fixtures
```bash
go run . --installFixturesDb=true --noServer=true
```
## create superuser
```bash
go run . --createSuperUser=true --noServer=true
```
## run and serve
```bash
go run .
```
feeling ready!
