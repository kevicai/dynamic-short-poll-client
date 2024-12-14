SERVER_BINARY = server/server
SERVER_PID_FILE = server.pid

.PHONY: server test run-server run-test clean

run-server:
	@echo "Starting server..."
	@nohup ./$(SERVER_BINARY) > server.log 2>&1 & echo $$! > $(SERVER_PID_FILE)
	@echo "Server started with PID $$(cat $(SERVER_PID_FILE))"

run-test:
	@echo "Running tests..."
	@go clean -testcache
	@go test ./tests -v

test: run-server run-test
	@echo "Tests completed."
	@make stop-server

stop-server:
	@if [ -f $(SERVER_PID_FILE) ]; then \
		kill -9 $$(cat $(SERVER_PID_FILE)) && echo "Stopped server with PID $$(cat $(SERVER_PID_FILE))"; \
		rm -f $(SERVER_PID_FILE); \
	else \
		echo "No server running."; \
	fi

clean:
	@echo "Cleaning up..."
	@make stop-server
	@rm -f server.log $(SERVER_PID_FILE)
