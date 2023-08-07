```
.
├── LICENSE
├── api
│   └── v1
│       ├── proto
│       │   └── goSentinel.proto
│       └── rpc
│           └── goSentinel.go
├── cmd
│   ├── app
│   │   ├── app.go
│   │   └── service.go
│   └── main.go
├── config.yaml
├── docker-compose.yaml
├── dockerfile
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   └── notifications
│       └── email
│           ├── config
│           │   └── config.go
│           ├── mail.go
│           └── model
│               └── interface.go
├── pkg
│   ├── db
│   │   ├── config
│   │   │   └── config.go
│   │   ├── db.go
│   │   └── helper
│   │       ├── enum_logical_operator_gen.go
│   │       └── filter.go
│   └── logger
│       ├── config
│       │   └── config.go
│       ├── logger.go
│       └── model
│           ├── interface.go
│           └── level.go
└── readme.md

20 directories, 24 files
```