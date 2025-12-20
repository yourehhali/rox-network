module seeder

go 1.24.0

replace social-network/pkg => ../pkg

require (
	github.com/mattn/go-sqlite3 v1.14.32
	social-network/pkg v0.0.0-00010101000000-000000000000
)

require golang.org/x/crypto v0.46.0 // indirect
