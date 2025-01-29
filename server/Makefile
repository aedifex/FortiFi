DB_NAME=FortiFi
DB_USER=root
SQL_FILE=init.sql

# Reset Database
clean:
	@echo "Dropping database $(DB_NAME)..."
	@echo "Requires mysql password"
	mysql -u$(DB_USER) -p -e "DROP DATABASE IF EXISTS $(DB_NAME);"
	@echo "Database $(DB_NAME) dropped."

# Run development server config
local-dev:
	go mod tidy
	go run cmd/server/main.go config=dev

# Setup database for development testing
database-dev:
	@echo "Creating database and tables..."
	@echo "Requires mysql password"
	mysql -u$(DB_USER) -p -e "SOURCE $(SQL_FILE);"
	@echo "Database and tables created successfully."