.PHONY: client go-build-client


client:
	cd ./client && ./client


go-build-client:
	go build -o ./client/client ./client