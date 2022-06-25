# go gin mongo studies

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-no-red.svg)](https://bitbucket.org/lbesson/ansi-colors)
[![Ask Me Anything !](https://img.shields.io/badge/Ask%20me-anything-1abc9c.svg)](https://GitHub.com/Naereen/ama)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)
[![Visual Studio Code](https://img.shields.io/badge/--007ACC?logo=visual%20studio%20code&logoColor=ffffff)](https://code.visualstudio.com/)
[![Docker](https://badgen.net/badge/icon/docker?icon=docker&label)](https://https://docker.com/)
[![Open Source? Yes!](https://badgen.net/badge/Open%20Source%20%3F/Yes%21/blue?icon=github)](https://github.com/Naereen/badges/)

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
