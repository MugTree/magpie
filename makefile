drop-data:
	rm magpie.db
	rm magpie.db-*

seed-data:
	touch magpie.db
	goose status
	goose up
	go run ./cmd/seed/main.go --urls=./cmd/seed/seed.csv --db=./magpie.db

lint:
	golangci-lint run .

dev: 
	air



debug:
	go build -gcflags="all=-N -l" -o ./tmp/server .

build-css: tailwindcss
	./tailwindcss -i tailwind.css -o app/public/css/app.css --minify

TAILWINDCSS_OS_ARCH := macos-arm64

tailwindcss:
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$(TAILWINDCSS_OS_ARCH)
	mv tailwindcss-$(TAILWINDCSS_OS_ARCH) tailwindcss
	chmod a+x tailwindcss
	mkdir -p node_modules/tailwindcss/lib && ln -sf tailwindcss node_modules/tailwindcss/lib/cli.js
	echo '{"devDependencies": {"tailwindcss": "latest"}}' >package.json

watch-css: tailwindcss
	./tailwindcss -i tailwind.css -o app/public/css/app.css --watch
