# Environment Variables
export DB_NAME=FortiFi
export DB_USER=root
export DB_PASS=""
export DB_URL=localhost:3306
export SIGNING_KEY=b2e138d8553ea7d7ff8731e87e41406277bd4c98
export FCM_KEY_PATH=/harness/firebase_auth.json
export PORT=:3000
export SQL_FILE=init.sql

.PHONY: clean local-dev database-dev clean-logs set-env

# Set environment variables explicitly
set-env:
	@echo "Setting environment variables..."
	@export DB_NAME=$(DB_NAME)
	@export DB_USER=$(DB_USER)
	@export DB_PASS=$(DB_PASS)
	@export DB_URL=$(DB_URL)
	@export SIGNING_KEY=$(SIGNING_KEY)
	@export FCM_KEY_PATH=$(FCM_KEY_PATH)
	@export PORT=$(PORT)
	@echo "Environment variables set."

# clean all
clean: clean-logs clean-database

# delete database
clean-database:
	@echo "Dropping database $(DB_NAME)..."
	@echo "Requires mysql password"
	mysql -u$(DB_USER) -p$(DB_PASS) -e "DROP DATABASE IF EXISTS $(DB_NAME);" || mysql -h 127.0.0.1 -P 3306 -u$(DB_USER) -p$(DB_PASS) -e "DROP DATABASE IF EXISTS $(DB_NAME);"
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
	mysql -u$(DB_USER) -p$(DB_PASS) -e "SOURCE $(SQL_FILE);" || mysql -h 127.0.0.1 -P 3306 -u$(DB_USER) -p$(DB_PASS) -e "SOURCE $(SQL_FILE);"
	@echo "Database and tables created successfully."

# Run development server config
local-dev: set-env
	go mod tidy
	config=dev go run ./cmd/server --verbose

# run tests locally -- this will reset the database
test: reset
	@echo "erasing cached tests"
	go clean -testcache
	@echo "running unit tests"
	config=dev go test ./cmd/server -v

cicd-database-dev:
	mysql -h 127.0.0.1 -P 3306 -u$(DB_USER) --password="" -e "SOURCE $(SQL_FILE);"
