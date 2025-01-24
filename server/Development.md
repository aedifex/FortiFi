/server
│
├── cmd/
│   └── server/
│       └── main.go -- server initialization and run
│
├── internal/
│   ├── handler/
│   │   └── user_handler.go -- routes handling
│   ├── database/
│   │   └── database.go -- database calls, structs, logic
│   └── middleware/
│       └── auth.go -- auth middleware
│
├── pkg/ -- reusable structs (loggers etc)
│   └── utils/
│       └── logger.go
│
├── api/
│   └── user_api.proto
│
├── configs/
│   └── config.yaml
│
├── web/
│   └── index.html
│
├── scripts/
│   └── setup.sh
│
├── go.mod
├── go.sum
├── README.md