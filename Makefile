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


# KIND cluster setup - creates the KIND cluster and local docker registry for use in this exercise
# https://kind.sigs.k8s.io/docs/user/quick-start/ ; https://kind.sigs.k8s.io/docs/user/local-registry/
# requires docker install as a prerequisite
create-cluster:
	@cd scripts; ./create-kind-cluster-registry.sh

# use if you have not done this: https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user
create-cluster-sudo:
	@cd scripts; sudo ./create-kind-cluster-registry.sh

test-cluster:
	kubectl cluster-info;kubectl get nodes

remove-cluster:
	kind cluster delete kind

remove-cluster-sudo:
	sudo kind cluster delete kind