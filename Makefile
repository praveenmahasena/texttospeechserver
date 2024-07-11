all: test
all: vet
all: package
all: package_race
all: run


test: vet
test: base_test

base_test:
	go test ./... -v

vet:
	go vet ./...

package: server

package_race: server_race

server:
	go build -o ./bin/server .

server_race:
	go build --race -o ./bin/server_race .

run:
	./bin/server
