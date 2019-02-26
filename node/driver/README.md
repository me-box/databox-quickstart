## Writing a driver in node

This driver will write a value to the key value store and then access it when the page is refreshed.  To get started, run:

```
cd src
npm install
npm start
```

This will start the test environment and start your code outside of databox for testing.

Then go to http://127.0.0.1:8080.  In the input box type some text and hit update.  Now refresh the page and you should see a statement "current config is: [your text]".  This means the driver has successfully set up a store to read/write to.

## Running on databox in dev mode

Databox supports starting a development version of your container with the source code mounted from you hosts filesystem. This is particularly usefully for complex apps that use multiple drivers, as they are hard to test externally. To do this we use the --devmount option of the databox command. The databoxDevSrcMnt variable at the top of the package.json holds the configuration json string for this.

>> **You will need to correct the path in databoxDevSrcMnt in package.json before this will work (set it to pwd).**

```
npm run build-dev           # Builds a dev image using the Dockerfile-dev (adds nodemon and npm_modules outside of the src path)
npm run start-databox-dev   # Starts a local copy of databox with --devmount set to point to the code here also setts the password to databoxDev

# wait for databox to start go to http://127.0.0.1 install the https certificate and then go to https://127.0.0.1
```

Finally, you'll need to upload your manifest file to tell databox about the new driver.

```
npm run upload-manifest     # Adds the databox manifest for this driver to databox
```

After this go to https://127.0.0.1 login and navigate to the app store where you should be able to see the driver ready to install.
Once installed you can edit the code on your host and changes should be visible to the running databox driver and nodemon will restart as required.

Databox maintains state between restarts so uploading the manifest and reinstalling is not always necessary between restarts.
If you make changes to the manifest this must be reuploaded and the driver reinstalled.


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
