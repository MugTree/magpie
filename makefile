seed-data:
	go run ./db/seed/generate.go --urls=./db/seed/seed.csv --db=./magpie.db

drop-data:
	rm magpie.db
	rm magpie.db-*
	rm ./markdown/*

recreate-db:
	touch magpie.db
	goose status
	goose up
	make seed-data 

lint:
	golangci-lint run .
