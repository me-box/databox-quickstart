NAME=go-test-driver
DEFAULT_REG=databoxsystems
VERSION=0.5.2

DATABOX_PASSWORD=databoxDev
DATABOX_DEV_SRC_MNT="[{\"ContName\":\"$(NAME)\",\"ContSrcPath\":\"/src\",\"HostSrcPath\":\"$(shell pwd)\"}]"

.PHONY: test
test:
	echo $(DATABOX_DEV_SRC_MNT)

.PHONY: all
all: build-amd64 build-arm64v8 publish-images

.PHONY: start
start:
	#Run the code outside of databox with a store and arbiter
	../../testenv/start.sh
	go run ./src/driver.go


.PHONY: restart-dev
restart-dev:
	#used in Dockerfile-dev to do rebuild and reload
	killall databox-dev-bin &
	go build -o ./bin/databox-dev-bin ./src/*.go
	./bin/databox-dev-bin

MATCH="\.go$$"
.PHONY: start-dev-mon
start-dev-mon:
	~/go/bin/reflex -r $(MATCH) -s make restart-dev

.PHONY: build-amd64
build-amd64:
	docker build -t $(DEFAULT_REG)/$(NAME)-amd64:$(VERSION) . $(OPTS)

.PHONY: build-amd64-dev
build-amd64-dev:
	docker build -t $(DEFAULT_REG)/$(NAME)-amd64:$(VERSION) -f Dockerfile-dev . $(OPTS)

.PHONY: build-arm64v8
build-arm64v8:
	docker build -t $(DEFAULT_REG)/$(NAME)-arm64v8:$(VERSION) -f Dockerfile-arm64v8 .  $(OPTS)


.PHONY: publish-images
publish-images:
	docker push $(DEFAULT_REG)/$(IMAGE_NAME)-amd64:$(VERSION)
	docker push $(DEFAULT_REG)/$(IMAGE_NAME)-arm64v8:$(VERSION)

.PHONY: start-databox
start-databox:
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock --network host -t $(DEFAULT_REG)/databox-amd64:$(VERSION) /databox start -sslHostName `hostname` --devmount=$(DATABOX_DEV_SRC_MNT) --release=$(VERSION) --password=$(DATABOX_PASSWORD)

.PHONY: stop-databox
stop-databox:
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock --network host -t $(DEFAULT_REG)/databox-amd64:$(VERSION) /databox stop

.PHONY: wipe-databox
wipe-databox:
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock --network host -t $(DEFAULT_REG)/databox-amd64:$(VERSION) /databox wipe -y


.PHONY: upload-manifest
upload-manifest:
	curl -sL --cookie-jar /tmp/databox-jar -H "authorization:Token $(DATABOX_PASSWORD)" --insecure --http1.1 -w "%{http_code}\n" 'https://127.0.0.1/api/connect' -o /dev/null
	curl -s  --insecure --http1.1 --cookie /tmp/databox-jar -H 'Content-Type: application/json' -X POST -d @databox-manifest.json -L 'https://127.0.0.1/app-store/ui/api/addManifest'

.PHONY: test
test:
#NOT IMPLIMENTED