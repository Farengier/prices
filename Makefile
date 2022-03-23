build:
	@go build -o build/db src/cli/database.go

run_db: build
	@build/db -c example/config.yml