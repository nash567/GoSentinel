```
.
├── LICENSE
├── Readme.md
├── api
│   └── v1
│       ├── pb
│       │   └── goSentinel
│       │       ├── goSentinel.pb.go
│       │       ├── goSentinel.pb.gw.go
│       │       └── goSentinel_grpc.pb.go
│       ├── proto
│       │   ├── buf.gen.yaml
│       │   ├── buf.lock
│       │   ├── buf.yaml
│       │   └── goSentinel.proto
│       └── rpc
│           └── goSentinel.go
├── cmd
│   ├── app
│   │   ├── app.go
│   │   └── service.go
│   └── main.go
├── config.yaml
├── diagrams
│   ├── application_registration.svg
│   ├── application_secret.svg
│   └── register_user.svg
├── docker-compose.yaml
├── dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── notifications
│   │   └── email
│   │       ├── config
│   │       │   └── config.go
│   │       ├── mail.go
│   │       └── model
│   │           └── interface.go
│   ├── resources
│   │   └── template
│   │       └── verification.html
│   └── service
│       ├── application
│       │   ├── config
│       │   │   └── config.go
│       │   ├── model
│       │   │   ├── filter.go
│       │   │   ├── interface.go
│       │   │   └── model.go
│       │   ├── repo
│       │   │   ├── filter.go
│       │   │   ├── query.go
│       │   │   └── repo.go
│       │   └── service.go
│       ├── auth
│       │   ├── auth.go
│       │   ├── config
│       │   │   └── config.go
│       │   ├── interceptor.go
│       │   ├── model
│       │   │   ├── interface.go
│       │   │   └── model.go
│       │   ├── repo
│       │   │   ├── query.go
│       │   │   └── repo.go
│       │   └── service.go
│       └── user
│           ├── config
│           ├── model
│           │   ├── interface.go
│           │   └── model.go
│           ├── repo
│           │   ├── filter.go
│           │   ├── query.go
│           │   └── repo.go
│           └── service.go
├── pkg
│   ├── cache
│   │   ├── interface.go
│   │   └── redis
│   │       ├── config
│   │       │   └── config.go
│   │       └── redis.go
│   ├── db
│   │   ├── config
│   │   │   └── config.go
│   │   ├── db.go
│   │   └── helper
│   │       ├── enum_logical_operator_gen.go
│   │       └── filter.go
│   └── logger
│       ├── config
│       │   └── config.go
│       ├── logger.go
│       └── model
│           ├── interface.go
│           └── level.go
└── update.txt

41 directories, 60 files
```


- go install google.golang.org/protobuf/cmd/protoc-gen-go@latest


## Application registration flow
![application registration](./diagrams/application_registration.svg)


## Application Credential Flow
![application credential flow](./diagrams/application_secret.svg)

## Register User Flow!
![register user flow](./diagrams/register_user.svg)