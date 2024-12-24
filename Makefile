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

test-v-ggtfs-common:
	go test --count=1 -v -tags=ggtfs_tests_common ./internal/pkg/ggtfs

coverage:
	JOURNEYS_BASE_URL=http://localhost:5678 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go test -count=1 -tags=all_tests -coverprofile cover.out ./...

coverage-html:
	JOURNEYS_BASE_URL=http://localhost:5678 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go test -count=1 -tags=all_tests -coverprofile cover.out ./... && go tool cover -html=cover.out -o coverage.html
tre:
	MEMCACHED_URL="localhost:11211" JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_BASE_URL="https://data.itsfactory.fi/journeys/api/1" JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go
tre-dry:
	MEMCACHED_URL="localhost:11211" JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_BASE_URL="https://data.itsfactory.fi/journeys/api/1" JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go --dry-run
tre-dev:
	JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_BASE_URL=http://localhost:5678 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go
vet:
	go vet ./...