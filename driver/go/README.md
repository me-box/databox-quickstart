## Writing an driver in go

go version 1.11 is recommended but 1.10 will work if the code for lib-go-databox installed in you go path.

If go version is less than 1.11 then:
* Install libzmq5
* Install libzmq3-dev
* run "go get github.com/gorrila.mux
* run "go get github.com/me-box/lib-go-databox

## Testing out side of databox

```
../../testenv/start.sh
make start
```
the ui will be available at http://127.0.0.1:8080/go-test-driver/ui

## Running on databox

To get running on the databox, you will first need to create a docker container to run your code.  To do this, you can either build the container on your databox, so pull your code, then in the src directory type:

For x86 platforms:
```
make build-amd64 DEFAULT_REG=databoxsystems VERSION=[running databox version eg 0.5.0 or latest]
```

For Arm v8 platforms:
```
make build-arm64v8 DEFAULT_REG=databoxsystems VERSION=[running databox version eg 0.5.0 or latest]
```

This will build and tag a docker image for use with databox. If databox is running on a machine other than the one you used to build the image then you will need to push it to docker hub under your own account. Change DEFAULT_REG to your docker hub registry, push it then pull it on the the target box the retag to databoxsystems.

Finally, you'll need to upload your manifest file to tell databox about the new app.  Log in to the databox and navigate to My Apps, then click on the "app store" app.  At the bottom of the page, use the form to upload your manifest.  Once uploaded, you can navigate to "App Store" and you should see go-test-app ready to install.
