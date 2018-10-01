## Databox quick start

This repository is for developers wishing to wriet drivers and apps for the databox platform.  For more information on databox, see our [project page](http://www.databoxproject.uk/) or [main github repo](https://github.com/me-box/databox)

##TLDR;

This repo contains all you need to create and test apps and drivers outside of the databox platform.  The only dependency that it has is an installation of Docker.  Note that this guide assumes that your base platform is MacOS or a flavour of linux - we do not current support development on Windows.

To simplify the development workflow, we have created a script that will set up a test environment to run your apps and drivers against. This allows you to test your code without having to build it as a docker container and install it on the platform.  Once you have a version you are happy with you can test it on the databox platform by following the instructions [here][Running on the databox platform] 

To start the test environment, run 
```
cd [rootdir]/databox-quickstart/testenv
chmod a+x start.sh
./start.sh
```

This will create two docker containers, zest and arbiter, which your testcode will communicate with to emulate communication with the databox.   To run a basic hello world example, go to the app/driver directory and choose your favoured language we currently support nodejs and golang, though python support is also in the pipeline.    Each of the directories contain the instructions you need to compile and run a basic "hello world" app/driver.

##Writing a driver

A databox driver is responsible for writing data to a datastore to make it available for apps.  Drivers are privileged code, so have unrestricted access to external addresses/ports.  We assume, in most cases that a databox will not have a externally accessible static IP address, so the typical approach is to require that drivers initiate communication to a datasource to gather data (rather than assuming the driver exposes a defined endpoint for datasources to connect to).  

Drivers (like apps) are written as webapps; user interaction with a driver is through a web interface; the top-level route for a driver is /ui, i.e. on the databox, when a user selects a driver, its webpage will be served from /ui.  Any additional rest endpoints must be under /ui (e.g. /ui/getConfig, ui/setConfig).  Drivers also have a manifest file, which must be written before it can be installed on a databox, but can be ignored for testing purposes. Details of the format of the manifest can be found [here][Running on the databox platform].  Each of the sample apps in this repo have a databox-manifest.json file which will probably be sufficient to get an idea of what is required.

Drivers will typically contain a set of initialisation steps:

1. Set up (one or more, store-dependent) clients to connect to its associated stores.
2. Request read/write access to the stores (which we assume have been previously authorised when the details of the driver manifest file are presented to the user at install time).

They will also typically contain logic that will:

3.  Provide a configuration endpoint for the user to provide an inital setup (for example, provide the IP address of an IoT Device or credentials for accessing a web-service).
4.  Provide the logic for connecting to a service or device OR provide an endpoint for actuating a device.
5.  Read/Write data to its stores. 

Steps 1 & 2 & 5 are most easily accomplished though use of the databox libraries.  It is possible to write components without the using the libraries, though this will require use a socket for writing binary data that conforms to our proprietary binary zeromq protocol (details can be found [here](https://me-box.github.io/zestdb/)).

##Writing an app

The process for writing an app is not to dissimilar from writing a driver, however, apps are untrusted code, and are therefore more restricted.  In particular they can only communicate with stores, and they are restricted to opening a port on 8080, to provide a web interface. Apps CANNOT directly access external addresses.  If they wish to do so, they must be defined explicitly in the manifest and must use the databox's export service.

##Testing on the databox platform

manifest,
tagging etc.
