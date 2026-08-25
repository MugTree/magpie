drop-data:
	rm magpie.db
	rm magpie.db-*
	rm ./markdown/*

seed-data:
	touch magpie.db
	goose status
	goose up
	go run ./db/seed/generate.go --urls=./db/seed/seed.csv --db=./magpie.db

lint:
	golangci-lint run .
