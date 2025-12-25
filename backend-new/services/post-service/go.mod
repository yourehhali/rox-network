module post-service

go 1.24.0

replace social-network/pkg => ../../pkg

require social-network/pkg v0.0.0-00010101000000-000000000000

require (
	github.com/golang-migrate/migrate v3.5.4+incompatible // indirect
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
	golang.org/x/crypto v0.46.0 // indirect
)
