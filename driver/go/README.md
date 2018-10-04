## Writing an driver in go

golang version 1.11 is recommended but 1.10 will work if lib-go-databox installed in you go path.

>> If go version is less than 1.11
>> Install libzmq5
>> Install libzmq3-dev
>> run "go get github.com/gorrila.mux
>> run "go get github.com/me-box/lib-go-databox

## Testing out side of databox

```
../../testenv/start.sh
make start
```

## Building an image to test on databox

```
make build-amd64 DEFAULT_REG=databoxsystems VERSION=0.5.0
```
