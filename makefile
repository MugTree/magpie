seed-db:
	go run ./db/seed/generate.go --urls=./db/seed/seed.csv --db=./magpie.db

drop-db:
	rm magpie.db
	rm magpie.db-*

recreate-db:
	touch magpie.db
	goose status
	goose up
	make seed-db 

lint:
	golangci-lint run .
