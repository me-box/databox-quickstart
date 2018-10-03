## Writing an app in node
This app will write a message to the helloworld actuator, created in the node driver code ([home]/databox-data-tracker/driver/node/src).  It also listens on new actuations, and when it observes one, it sends it to the client over a websocket. To get started, first ensure that your are running the test environment (i.e. you have called [home]/databox-quickstart/testenv/start.sh) and that you are running the driver code.  Then do:

```
cd src
npm install
npm start
```

Then go to http://127.0.0.1:8090/ui/actuate.  You should see the following in the terminal logs:

```
[Creating TEST http server] 8090
started listening to actuator
seen data from the hello world actuator! { msg: '1538476623699:test actuation!' }
successfully actuated!
```

and in the driver terminal:

```
[Actuation] data received  { timestamp: 1538476623730,
  datasourceid: undefined,
  key: '',
  data: '{"msg":"1538476623699:test actuation!"}' }
```

## Running on databox

To get running on the databox, you will first need to create a docker container to run your code.  To do this, you can either build the container on your databox, so pull your code, then in the src directory type:

```
npm run docker
```

Or you could push the container image to a docker repository (e.g dockerhub), then pull it onto the databox.  So do the following:

```
npm run docker
docker tag [dockerhubusername]/databox-app-helloworld-node
docker push [dockerhubusername]/databox-app-helloworld-node
```

Once your image is on the databox, you need to name it so that it can be found by the container manager.  By default the container manager looks for containers in databoxsystems.  Assuming you are running version 0.5.1 on an Intel 64 bit machine, do the following to rename your image:

```
docker tag [dockerhubusername]/databox-app-helloworld-node databoxsystems/databox-app-helloworld-node-amd64:0.5.1
```

or if you are running the latest bleeding-edge version of the platform you can do:

```
docker tag [dockerhubusername]/databox-app-helloworld-node databoxsystems/databox-app-helloworld-node-amd64
```

Finally, you'll need to upload your manifest file to tell databox about the new app.  Log in to the databox and navigate to My Apps, then click on the "app store" app.  At the bottom of the page, use the form to upload your manifest.  Once uploaded, you can navigate to "App Store" and you should see databox-app-helloworld-node ready to install. 