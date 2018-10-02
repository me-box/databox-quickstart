##Writing a driver in node

This driver will write a value to the key value store and then access it when the page is refreshed.  To get started, first ensure that your are running the test environment (i.e. you have called [home]/databox-quickstart/testenv/start.sh). Then run:

```
cd src
npm install
npm run testmode
```

Then go to http://127.0.0.1:8080.  In the input box type some text and hit update.  Now refresh the page and you should see a statement "current config is: [your text]".  This means the driver has successfully set up a store to read/write to.    

##Running on databox

To get running on the databox, you will first need to create a docker container to run your code.  To do this, you can either build the container on your databox, so pull your code, then in the src directory type:

```
npm run docker
```

Or you could push the container image to a docker repository (e.g dockerhub), then pull it onto the databox.  So do the following:

```
npm run docker
docker tag [dockerhubusername]/databox-driver-helloworld-node
docker push [dockerhubusername]/databox-driver-helloworld-node
```

Once your image is on the databox, you need to name it so that it can be found by the container manager.  By default the container manager looks for containers in databoxsystems.  Assuming you are running version 0.5.1 on an Intel 64 bit machine, do the following to rename your image:

```
docker tag [dockerhubusername]/databox-driver-helloworld-node databoxsystems/databox-driver-helloworld-node-amd64:0.5.1
```

or if you are running the latest bleeding-edge version of the platform you can do:

```
docker tag [dockerhubusername]/databox-driver-helloworld-node databoxsystems/databox-driver-helloworld-node-amd64
```

Finally, you'll need to upload your manifest file to tell databox about the new app/driver.  Log in to the databox and navigate to My Apps, then click on the "app store" app.  At the bottom of the page, use the form to upload your manifest.  Once uploaded, you can navigate to "App Store" and you should see it ready to install. 