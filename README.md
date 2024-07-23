# websocketCollab
Collaborative system for editing text backend made on Go. This project was made to learn Fiber framework and websockets.

## Stack
I used following technologies:
- Go as programming language
- Fiber as web framework
- SQLite (or Postgres in Docker) as DB
- Websockets
- Docker

## Run project
You can run project with docker or without it

### Running project in Docker
To run project in Docker:
1) Uncomment PostgreSQL connection string in cmd/websocketCollabServer/main.go, and comment SQLite connection string:
```
// this connection can be used to connect to postgres inside docker container.
// db, err := gorm.Open(postgres.Open(getDsn()), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

// this connection uses sqlite
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
```

2) Build image:
`docker build . -t "wscollab:latest"`
3) Run docker-compose:
`docker-compose -f docker-compose.yaml up`

### Running project without Docker
To run project without docker with SQLite Database
1) Use command `go mod download` to install dependencies
2) Use command `go run cmd/websocketCollabServer/main.go` to run project

## How to use it?
API of this project includes 2 endpoints:
1) POST - `/api/v1/documents` - will create new empty document. This request requires no parameters or body, in response you will get JSON-encoded document, which will include Id
2) GET - `/api/v1/documents/{id}` - websocket handshake path, connect by id of your document using any websocket client, for example https://chromewebstore.google.com/detail/simple-websocket-client/pfdhoblngboilpfeibdedpjgfnlcodoo?pli=1

### How to add text?
To add text, you need to send json-encoded message with this format:
```
{
  "cmd": {command type here, command can be "INSERT" or "DELETE"},
  "position": {float32 here, document CRDT uses float indexing to avoid conflicts},
  "value": {byte value here, if command is "DELETE" value will be ignored}
}
```
After sending such message, all clients, connnected to this document will recieve edited document text as a message
