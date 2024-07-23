build :
	go build -o ./bin/go-blockchain
run: build
	./bin/go-blockchain
clean:
	go clean
test:
	go test ./...