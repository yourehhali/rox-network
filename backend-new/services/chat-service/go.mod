module chat-service

go 1.24.2

replace social-network/pkg => ../../pkg

require (
	github.com/gorilla/websocket v1.5.3
	social-network/pkg v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang-migrate/migrate v3.5.4+incompatible // indirect
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
)
