#!/bin/sh
echo "Starting databox test env"

docker kill zest
docker kill arbiter

ZEST_IMAGE_VERSION="jptmoore/zest:v0.1.0"
ARBITER_IMAGE_VERSION="jptmoore/arbiter:latest"


echo "start the arbiter"
docker run -p 4444:4444 -p 4445:4445 -d --name arbiter --rm ${ARBITER_IMAGE_VERSION} /app/zest/server.exe --request-endpoint tcp://0.0.0.0:4444 --secret-key-file example-server-key --token-key-file example-token-key --enable-logging

echo "register a store"
docker run --network host -it ${ZEST_IMAGE_VERSION} /app/zest/client.exe --server-key 'vl6wu0A@XP?}Or/&BR#LSxn>A+}L)p44/W[wXL3<' --request-endpoint tcp://0.0.0.0:4444 --path '/cm/upsert-container-info' --mode post --payload "{\"name\": \"127.0.0.1\", \"type\": \"store\", \"key\": \"storeSecret\"}" --token secret

echo "Get the store token"
docker run --network host -it ${ZEST_IMAGE_VERSION} /app/zest/client.exe --server-key 'vl6wu0A@XP?}Or/&BR#LSxn>A+}L)p44/W[wXL3<' --request-endpoint tcp://0.0.0.0:4444 --path '/store/secret' --mode get --identity '127.0.0.1' --token storeSecret > /tmp/storekeyextra.txt
# strip
tr -d '\r\n' < /tmp/storekeyextra.txt > /tmp/storekey.txt
echo "store token is "
cat /tmp/storekey.txt

echo "start the store with the key from above"
docker run -p 5555:5555 -p 5556:5556 -d --name zest -v /tmp/storekey.txt:/storekey.txt --rm ${ZEST_IMAGE_VERSION} /app/zest/server.exe --secret-key-file example-server-key --token-key-file /storekey.txt --identity '127.0.0.1' --enable-logging

HOST_NAME=$(hostname)
echo "registering app called " $HOST_NAME
docker run --network host -it ${ZEST_IMAGE_VERSION} /app/zest/client.exe --server-key 'vl6wu0A@XP?}Or/&BR#LSxn>A+}L)p44/W[wXL3<' --request-endpoint tcp://0.0.0.0:4444 --path '/cm/upsert-container-info' --mode post --payload "{\"name\": \"${HOST_NAME}\", \"type\": \"app\", \"key\": \"secret\"}" --token secret

echo "Granting permissions for app called " $HOST_NAME
docker run --network host -it ${ZEST_IMAGE_VERSION} /app/zest/client.exe --server-key 'vl6wu0A@XP?}Or/&BR#LSxn>A+}L)p44/W[wXL3<' --path '/cm/grant-container-permissions' --mode post --payload "{\"name\": \"${HOST_NAME}\", \"caveats\": [], \"route\": {\"method\": \"POST\", \"path\": \"/*\", \"target\": \"127.0.0.1\"}}" --token secret --request-endpoint 'tcp://127.0.0.1:4444' --token secret
docker run --network host -it ${ZEST_IMAGE_VERSION} /app/zest/client.exe --server-key 'vl6wu0A@XP?}Or/&BR#LSxn>A+}L)p44/W[wXL3<' --path '/cm/grant-container-permissions' --mode post --payload "{\"name\": \"${HOST_NAME}\", \"caveats\": [], \"route\": {\"method\": \"GET\", \"path\": \"/*\", \"target\": \"127.0.0.1\"}}" --token secret --request-endpoint 'tcp://127.0.0.1:4444' --token secret
docker run --network host -it ${ZEST_IMAGE_VERSION} /app/zest/client.exe --server-key 'vl6wu0A@XP?}Or/&BR#LSxn>A+}L)p44/W[wXL3<' --path '/cm/grant-container-permissions' --mode post --payload "{\"name\": \"${HOST_NAME}\", \"caveats\": [], \"route\": {\"method\": \"DELETE\", \"path\": \"/*\", \"target\": \"127.0.0.1\"}}" --token secret --request-endpoint 'tcp://127.0.0.1:4444' --token secret
