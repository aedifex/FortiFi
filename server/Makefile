DB_NAME=FortiFi
DB_USER=root
SQL_FILE=init.sql

.PHONY: clean local-dev database-dev clean-logs

# clean all
clean: clean-logs clean-database

# Reset Database
clean-database:
	@echo "Dropping database $(DB_NAME)..."
	@echo "Requires mysql password"
	mysql -u$(DB_USER) -p -e "DROP DATABASE IF EXISTS $(DB_NAME);"
	@echo "Database $(DB_NAME) dropped."

# Reset logs
clean-logs:
	@echo "Deleting log files"
	rm server.log*

# Run development server config
local-dev:
	go mod tidy
	config=dev go run ./cmd/server

# Setup database for development testing
database-dev:
	@echo "Creating database and tables..."
	@echo "Requires mysql password"
	mysql -u$(DB_USER) -p -e "SOURCE $(SQL_FILE);"
	@echo "Database and tables created successfully."