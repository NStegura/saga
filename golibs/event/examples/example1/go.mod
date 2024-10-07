module github.com/NStegura/saga/golibs/event/examples/example1

go 1.22.0

replace github.com/NStegura/saga/golibs/event => ../../

require (
	github.com/NStegura/saga/golibs/event v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v5 v5.7.1
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
