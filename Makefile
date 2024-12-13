SERVER_BINARY = server/server

.PHONY: server test run-test

run-server:
	@ echo "Starting server..."
	@ ./$(SERVER_BINARY)

run-test: 
	@ echo "Running tests..."
	@ go clean -testcache 
	@ go test ./tests -v

test: run-server run-test
	@echo "Tests completed."
	@pkill -f $(SERVER_BINARY) # stop the server after tests