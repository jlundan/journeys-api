.PHONY: build, test, vendor

mods:
	go mod download

vendor:
	go mod vendor

build:
	GOFLAGS="-mod=vendor" go build -o ./bin/journeys.api cmd/journeys/journeys.go

test:
	GOFLAGS="-mod=vendor" JOURNEYS_BASE_URL=http://localhost:5678 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go test --count=1 -tags=all_tests ./...

test-v:
	GOFLAGS="-mod=vendor" JOURNEYS_BASE_URL=http://localhost:5678 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go test --count=1 -v -tags=all_tests ./...

test-v-ggtfs:
	GOFLAGS="-mod=vendor" go test --count=1 -v -tags=ggtfs_tests ./internal/pkg/ggtfs

test-v-journeys:
	GOFLAGS="-mod=vendor" go test --count=1 -v -tags=journeys_tests ./internal/app/journeys/...

test-v-ggtfs-common:
	GOFLAGS="-mod=vendor" go test --count=1 -v -tags=ggtfs_tests_common ./internal/pkg/ggtfs

coverage:
	GOFLAGS="-mod=vendor" JOURNEYS_BASE_URL=http://localhost:5678 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go test -count=1 -tags=all_tests -coverprofile cover.out ./...

coverage-html:
	GOFLAGS="-mod=vendor" JOURNEYS_BASE_URL=http://localhost:5678 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go test -count=1 -tags=all_tests -coverprofile cover.out ./... && go tool cover -html=cover.out -o coverage.html
tre:
	GOFLAGS="-mod=vendor" MEMCACHED_URL="localhost:11211" JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_BASE_URL="https://data.itsfactory.fi/journeys/api/1" JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go start
tre-custom-cache:
	GOFLAGS="-mod=vendor" MEMCACHED_URL="localhost:11211" JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_BASE_URL="https://data.itsfactory.fi/journeys/api/1" JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" JOURNEYS_SHORT_CACHE_LOWER_BOUND="0" JOURNEYS_SHORT_CACHE_UPPER_BOUND="5" JOURNEYS_SHORT_CACHE_DURATION="10s" JOURNEYS_LONG_CACHE_DURATION="30s" go run cmd/journeys/journeys.go start
tre-no-cache:
	GOFLAGS="-mod=vendor" MEMCACHED_URL="localhost:11211" JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_BASE_URL="https://data.itsfactory.fi/journeys/api/1" JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go start --disable-cache
tre-dry:
	GOFLAGS="-mod=vendor" MEMCACHED_URL="localhost:11211" JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_BASE_URL="https://data.itsfactory.fi/journeys/api/1" JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go start --dry-run
tre-dev:
	GOFLAGS="-mod=vendor" JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_PORT=5678 JOURNEYS_BASE_URL=http://localhost:5678/v1 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go start --disable-cache
tre-dev-no-val:
	GOFLAGS="-mod=vendor" JOURNEYS_GTFS_PATH=.gtfs JOURNEYS_PORT=5678 JOURNEYS_BASE_URL=http://localhost:5678/v1 JOURNEYS_VA_BASE_URL="https://data.itsfactory.fi/journeys/api/1" go run cmd/journeys/journeys.go start --disable-cache --skip-validation
vet:
	GOFLAGS="-mod=vendor" go vet ./...