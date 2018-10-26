IMAGE_NAME=go-test-app
DEFAULT_REG=databoxdev
VERSION=dev

.PHONY: all
all: build-amd64 build-arm64v8 publish-images

.PHONY: start
start:
	go run ./src/app.go

.PHONY: build-amd64
build-amd64:
	docker build -t $(DEFAULT_REG)/$(IMAGE_NAME)-amd64:$(VERSION) . $(OPTS)

.PHONY: build-arm64v8
build-arm64v8:
	docker build -t $(DEFAULT_REG)/$(IMAGE_NAME)-arm64v8:$(VERSION) -f Dockerfile-arm64v8 .  $(OPTS)

.PHONY: publish-images
publish-images:
	docker push $(DEFAULT_REG)/$(IMAGE_NAME)-amd64:$(VERSION)
	docker push $(DEFAULT_REG)/$(IMAGE_NAME)-arm64v8:$(VERSION)

.PHONY: test
test:
#NOT IMPLIMENTED