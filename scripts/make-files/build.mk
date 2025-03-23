run:
	clear
	mkdir -p bin
	go build -o ./bin/api ./cmd/api
	ENV=local CONFIGURATION_FILEPATH=$(PWD)/cmd/api/.env ./bin/api