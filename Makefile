SERVER_BINARY = server/server
SERVER_SOURCE = server/server.go

server:
	echo "Starting server..."
	go build -o $(SERVER_BINARY) $(SERVER_SOURCE)
	./$(SERVER_BINARY)

clean:
	rm -f $(SERVER_BINARY)