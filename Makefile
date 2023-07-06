.PHONY: all clean

BINARY_SERVER := server
BINARY_CLIENT := client

all: $(BINARY_SERVER) $(BINARY_CLIENT)

$(BINARY_SERVER):
	go build -o $(BINARY_SERVER) ./server.go

$(BINARY_CLIENT):
	go build -o $(BINARY_CLIENT) ./client.go

clean:
	rm -f $(BINARY_SERVER) $(BINARY_CLIENT)