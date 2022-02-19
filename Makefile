build:
	protoc --micro_out=. --go_out=. *.proto
	go mod tidy
