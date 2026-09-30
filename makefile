production-build-app:
	CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external -extldflags -static" -o ./bin/app.amd64 .

build-app:
	go build -o ./bin/app.mac .

drop-data:
	rm billiards.db
	rm billiards.db-*

seed-data:
	touch billiards.db
	goose status
	goose up
	go run ./app/db/seed/generate.go --urls=./app/db/seed/seed.csv --db=./billiards.db

lint:
	golangci-lint run .

run-tests:
	cd www && go test -v && cd ..

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