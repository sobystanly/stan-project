BINARY_NAME=risks

## build: Build binaries
build:
	@echo "Building..."
	@env CGO_ENABLED=0 go build -ldflags="-s -w" -o ${BINARY_NAME}
	@echo "Built!"

## run: builds and runs the application
run: build
	@echo "Starting..."
	@env ./${BINARY_NAME} &
	@echo "Started!"

## clean: runs go clean and delete binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned"

## start:  an alias to run
start: run

## stop: stops the running applications
stop:
	@echo "Stopping"
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the applications
restart: stop start

## test: runs all the tests
test:
	go test -v ./... # Run tests for the risks service