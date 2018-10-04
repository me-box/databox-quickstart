FROM node:alpine

RUN apk add --update make gcc g++ python curl git krb5-dev zeromq-dev && \
    npm install zeromq --zmq-external --save

ADD ./package.json package.json
RUN npm install

RUN apk del make gcc g++ python curl git krb5-dev

ADD ./main.js main.js
ADD ./views views

LABEL databox.type="app"
LABEL databox.manifestURL="/helloworld-node-driver/databox-manifest.json"
EXPOSE 8080

CMD ["npm","start"]
