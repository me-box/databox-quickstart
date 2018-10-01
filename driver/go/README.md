## Writing an driver in go

golang version 1.11 is recommended but 1.10 will work if lib-go-databox installed in you go path.

#Testing out side of databox

```
../../testenv/start.sh
make start
```

# Building an image to test on databox

```
make build
```
make build-amd64 DEFAULT_REG=databoxsystems VERSION=0.5.0
```
