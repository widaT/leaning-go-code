export GOCMD=GO111MODULE=on CGO_ENABLED=0 go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

all : service web cli

service :
	@echo "build service"
	@mkdir -p build
	$(GOBUILD) -a -installsuffix cgo -ldflags '-w' -o build/service cmd/service/*.go

web :
	@echo "build web"
	$(GOBUILD) -a -installsuffix cgo -ldflags '-w' -o build/web cmd/web/*.go

cli :
	@echo "build cli"
	$(GOBUILD) -a -installsuffix cgo -ldflags '-w' -o build/cli cmd/cli/*.go
	
docker:
	@echo "build docker images"
	docker-compose up --build
	
.PHONY: clean
clean:
	@rm -rf build/

.PHONY: proto
proto:
	protoc --go_out=. --micro_out=. test.proto





