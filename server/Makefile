DB_NAME=FortiFi
DB_USER=root
SQL_FILE=init.sql
#DB_PASSWORD=_BrIGChRf5lm
DB_PASSWORD = "Password123!"
.PHONY: clean local-dev database-dev clean-logs

# clean all
clean: clean-logs clean-database

setup-db:
	@echo "Resetting MySQL root password..."
	mysql --connect-expired-password -u root -p'$(MYSQL_ROOT_PASSWORD)' -e "ALTER USER 'root'@'localhost' IDENTIFIED BY '$(DB_PASSWORD)'; FLUSH PRIVILEGES;"

# delete database
clean-database:
	@echo "Dropping database $(DB_NAME)..."
	@echo "Requires mysql password"
	mysql --connect-expired-password -u$(DB_USER) -p$(DB_PASSWORD) -e "DROP DATABASE IF EXISTS $(DB_NAME);"
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
	mysql --connect-expired-password -u$(DB_USER) -p$(DB_PASSWORD) -e "SOURCE $(SQL_FILE);"
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

cicd-database-dev:
	mysql -h 127.0.0.1 -P 3306 -u$(DB_USER) --password="" -e "SOURCE $(SQL_FILE);"

production:
	go mod tidy
	config=prod go run ./cmd/server &
	disown
