DB_NAME=FortiFi
DB_USER=root
SQL_FILE=init.sql
DB_PASSWORD=${DB_PASSWORD}
.PHONY: clean local-dev database-dev clean-logs

# clean all
clean: clean-logs clean-database

# delete database
clean-database:
	@echo "Dropping database $(DB_NAME)..."
	@echo "Requires mysql password"
	mysql -u$(DB_USER) -p$(DB_PASSWORD) -e "DROP DATABASE IF EXISTS $(DB_NAME);"
	@echo "Database $(DB_NAME) dropped."

# delete logs
clean-logs:
	@echo "Deleting log files"
	@find . -type f -name "*.log*" -exec rm -f {} + && echo "Files deleted!" || echo "No log files found"

# reset database
reset: clean-database database-dev
	
# Setup database for development testing
database-dev:
	@echo "Creating database and tables..."
	@echo "Requires mysql password"
	mysql -u$(DB_USER) -p$(DB_PASSWORD) -e "SOURCE $(SQL_FILE);"
	@echo "Database and tables created successfully."

# Run development server config
local-dev:
	go mod tidy
	config=dev go run ./cmd/server --verbose

# run tests locally -- this will reset the database
test: reset
	@echo "erasing cached tests"
	go clean -testcache
	@echo "running unit tests"
	config=dev go test ./cmd/server -v