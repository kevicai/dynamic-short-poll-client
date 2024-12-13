SERVER_BINARY = server
SERVER_SOURCE = server.go

.PHONY: clean server

run:
	@ echo "Starting server..."
	@ go build -o $(SERVER_BINARY) $(SERVER_SOURCE)
	@ ./$(SERVER_BINARY)

clean:
	@ rm -f $(SERVER_BINARY)