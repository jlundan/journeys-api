.PHONY: build, test

mods:
	go mod download

build:
	go build -o ./bin/journeys.api cmd/journeys/journeys.go

test:
	JOURNEYS_BASE_URL=http://localhost:5678 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go test --count=1 -tags=all_tests ./...

test-v:
	JOURNEYS_BASE_URL=http://localhost:5678 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go test --count=1 -v -tags=all_tests ./...

test-v-ggtfs:
	go test --count=1 -v -tags=ggtfs_tests ./internal/pkg/ggtfs

test-v-journeys:
	go test --count=1 -v -tags=journeys_tests ./internal/app/journeys/...

test-v-ggtfs-common:
	go test --count=1 -v -tags=ggtfs_tests_common ./internal/pkg/ggtfs

coverage:
	JOURNEYS_BASE_URL=http://localhost:5678 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go test -count=1 -tags=all_tests -coverprofile cover.out ./...

coverage-html:
	JOURNEYS_BASE_URL=http://localhost:5678 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go test -count=1 -tags=all_tests -coverprofile cover.out ./... && go tool cover -html=cover.out -o coverage.html
tre:
	MEMCACHED_URL="localhost:11211" JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_BASE_URL="https://data.itsfactory.fi/journeys/api/1" JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go start
tre-no-cache:
	MEMCACHED_URL="localhost:11211" JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_BASE_URL="https://data.itsfactory.fi/journeys/api/1" JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go start --disable-cache
tre-dry:
	MEMCACHED_URL="localhost:11211" JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_BASE_URL="https://data.itsfactory.fi/journeys/api/1" JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go start --dry-run
tre-dev:
	JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_PORT=5678 JOURNEYS_BASE_URL=http://localhost:5678/v1 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go start --disable-cache
tre-dev-no-val:
	JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_PORT=5678 JOURNEYS_BASE_URL=http://localhost:5678/v1 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go start --disable-cache --skip-validation
vet:
	go vet ./...