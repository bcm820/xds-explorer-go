.DEFAULT_GOAL := build
 
#------------------------------------------------------------------------------
#-- Common
#------------------------------------------------------------------------------

.PHONY: build
build: vendor
	@echo "Building binary..."
	@go build --mod=vendor

.PHONY: vendor
test: vendor
	@echo "Testing..."
	@go test --mod=vendor

.PHONY: vendor
vendor:
	@echo "Vendoring dependencies..."
	@go mod vendor

#------------------------------------------------------------------------------
#-- Docker
#------------------------------------------------------------------------------

.PHONY: docker
docker: vendor
	@echo "Generating docker image..."
	@docker build -f docker/Dockerfile -t bcmendoza/xds-explorer:latest .