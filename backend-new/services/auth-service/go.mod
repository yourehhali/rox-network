module auth-service

go 1.24.0

replace social-network/pkg => ../../pkg

require social-network/pkg v0.0.0-00010101000000-000000000000

require (
	github.com/go-gorp/gorp/v3 v3.1.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
	github.com/rubenv/sql-migrate v1.8.1 // indirect
	golang.org/x/crypto v0.46.0 // indirect
)
