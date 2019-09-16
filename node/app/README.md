## Writing an app in node

This app will write a message to the helloworld actuator, created in the node driver code ([home]/databox-data-tracker/driver/node/src).  It also listens on new actuations, and when it observes one, it sends it to the client over a websocket. To get started, first ensure that your are running the test environment (i.e. you have called [home]/databox-quickstart/testenv/start.sh) and that you are running the driver code.  Then do:

```
cd src
npm install
npm start
```

This will start the test environment and start your code outside of databox for testing.

To do an end to end test you will also need to run the driver. In a separate terminal window run:

```
cd [to driver src path]
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

## Running on databox in dev mode

Databox supports starting a development version of your container with the source code mounted from you hosts filesystem. This is particularly usefully for complex apps that use multiple drivers, as they are hard to test externally. To do this we use the --devmount option of the databox command. The databoxDevSrcMnt variable at the top of the package.json holds the configuration json string for this.

>> **You will need to correct the path in databoxDevSrcMnt in package.json before this will work (set it to pwd).**

```
npm run build-dev           # Builds a dev image using the Dockerfile-dev (adds nodemon and npm_modules outside of the src path)
npm run start-databox-dev   # Starts a local copy of databox with --devmount set to point to the code here also setts the password to databoxDev

# wait for databox to start go to http://127.0.0.1 install the https certificate and then go to https://127.0.0.1
```

Finally, you'll need to upload your manifest file to tell databox about the new app.

```
npm run upload-manifest     # Adds the databox manifest for this app to databox
```

After this go to https://127.0.0.1 login and navigate to the app store where you should be able to see the app ready to install. If the driver is not installed you will be asked to install it first if its manifest has been uploaded.

Once installed you can edit the code on your host and changes should be visible to the running databox app and nodemon will restart as required.

Databox maintains state between restarts so uploading the manifest and reinstalling is not always necessary between restarts.
If you make changes to the manifest this must be reuploaded and the app reinstalled.


## Running on databox in production mode

To get running on the databox, you will first need to create a docker container to run your code.  To do this, in the src directory type:

```
npm run build-prod       # Builds a production image using the Dockerfile
npm run start-databox   # Starts a local copy of databox and sets the password to databoxDev

# wait for databox to start go to http://127.0.0.1 install the https certificate and then go to https://127.0.0.1

```

Finally, you'll need to upload your manifest file to tell databox about the new driver.

```
npm run upload-manifest     # Adds the databox manifest for this driver to databox
```

In this mode if you make changes to the code you must run `npm run build-prod` and restart the driver using the restart icon in the top left of the ui.
If you make changes to the manifest this must be reuploaded and the driver reinstalled.

# Stopping and resetting databox

To stop databox run:

```
npm run stop-databox
```

To completely reset databox run:

```
npm run wipe-databox
```

# Export Destination

When actuated the app also tries to post data to an external URL using
the export service. The default URL is `https://postman-echo.com/post`.
If you change the URL then make sure you change it in the manifest
export-whitelist, upload the changed manifest and re-install the app,
as well as changing the URL in the code (`main.js`).

