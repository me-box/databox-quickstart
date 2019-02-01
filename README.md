## Databox quick start

This repository is for developers wishing to write drivers and apps for the databox platform.  For more information on databox, see [what is databox](https://github.com/me-box/databox/blob/master/documents/what-is-databox.md) and our [project page](http://www.databoxproject.uk/) or [main github repo](https://github.com/me-box/databox)

## TLDR

This repo contains all you need to create and test apps and drivers outside of the databox platform.  The only dependency that it has is an installation of Docker.  Note that this guide assumes that your base platform is MacOS or a flavour of linux - we do not current support development on Windows. Instruction on setting up docker can be found [on the docker website](https://docs.docker.com/install/#supported-platforms)

To simplify the development workflow, we have created a script that will set up a test environment outside to run your apps and drivers against and a method to run them inside databox and mount code from the local host. This allows you to test your code without having to build it as a docker container and install it on the platform.  Once you have a version you are happy with you can test it on the databox platform by following the instructions [here](#testing-on-the-databox-platform)

To start the test environment, run
```
cd databox-quickstart/
./testenv/start.sh
```

To stop the test environment, run
```
cd databox-quickstart/
./testenv/stop.sh
```

This will create two docker containers, zest (a data store) and arbiter, which your test code will communicate with to emulate communication with the databox.  To run a basic hello world example, go to the app/driver directory with in then language of your choice, we currently support from nodejs or golang, though python support is also in the pipeline. Each of the directories contain the instructions you need to compile and run a basic "hello world" app/driver.

## Databox architecture crash course

Overview of the databox architecture:

![databox architecture](https://github.com/tlodge/databox-sdk-tutorial/blob/master/images/overview/databoxoverview.svg)

Databox composes of *data sources*, *data stores*, the *arbiter*, *drivers*, *apps* and *export service*.

* A data source represents data from either some physical hardware (such as a sensor, an IoT device in the home, e.g. Philips Hue Bulbs, smartplugs) or a cloud service (e.g. Twitter, Facebook, Gmail).

* A driver is a piece of software that is installed on the databox to communicate with a specific device or service to create a set of data sources.

* The arbiter is the keeper of permissions and the minter of tokens. The tokens minted by the arbiter can be independently verified and authenticated by the Data stores. Any component wishing to read or write data must present a valid token with the appropriate permissions.

* Data stores are access controlled and auditable database of data sources. They support structured and unstructured time-series data and a key-value store. To access data in a store you must have a valid token from the arbiter.

* Apps, are data processors they do not have direct access to drivers or the internet. All they know about is data sources they have access to. All apps have a manifest file, which sets out the data sources it will require access to. At install time, if the user accepts the details of the manifest, then the arbiter will set up the necessary permissions. Apps can also request permission to export data to an external service via the export service.

* The export service allows apps to send data out of databox. Data that passes though it is logged and auditable by the user.

All components in the databox run as docker containers, and it is the container managers job to pull, run and manage their life cycle containers. Communication between components is also tightly controlled by the databox network, which is beyond the scope of this guide.

## Writing a driver

A databox driver is responsible for writing data to a datastore to make it available for apps.  Drivers are privileged code, so have unrestricted access to the local network and can request external access in their manifests (using the ExternalWhitelist).  We assume, in most cases that a databox will not have a externally accessible static IP address, so the typical approach is to require that drivers initiate communication to a datasource to gather data (rather than assuming the driver exposes a defined endpoint for datasources to connect to).

Drivers (like apps) are written as webapps; user interaction with a driver is through a web interface; the top-level route for a driver is /ui, i.e. on the databox, when a user selects a driver, its webpage will be served from /ui.  Any additional rest endpoints must be under /ui (e.g. /ui/getConfig, /ui/setConfig).  Drivers also have a manifest file, which must be written before it can be installed on a databox, but can be ignored for testing purposes. Details of the format of the manifest can be found [here][Running on the databox platform].  Each of the sample apps in this repo have a databox-manifest.json file which will probably be sufficient to get an idea of what is required.

Drivers will typically contain a set of initialisation steps:

1. Set up (one or more, store-dependent) clients to connect to its associated stores.
2. Request read/write access to the stores (which we assume have been previously authorised when the details of the driver manifest file are presented to the user at install time).

They will also typically contain logic that will:

3.  Provide a configuration endpoint for the user to provide an initial setup (for example, provide the IP address of an IoT Device or credentials for accessing a web-service).
4.  Provide the logic for connecting to a service or device OR provide an endpoint for actuating a device.
5.  Read/Write data to its stores.

Steps 1 & 2 & 5 are most easily accomplished though use of the databox libraries.  It is possible to write components without the using the libraries, although this will require use of a socket for writing binary data that conforms to our proprietary binary zeromq protocol (details can be found [here](https://me-box.github.io/zestdb/)).

## Writing an app

The process for writing an app is not too dissimilar from writing a driver, however, apps are untrusted code, and are therefore more restricted.  In particular they can only communicate with stores, and they are restricted to opening a port on 8080, to provide a web interface. Apps CANNOT directly access external addresses or the local network.  If they wish to access external addresses apps must defined this explicitly in the manifest (in the ExportWhitelists) and must use the databox's export service.

## Testing on the databox platform

These instructions assume you have a working databox.  If you do not, please read the instructions [here](https://github.com/me-box/databox). To test your app/driver on the databox platform you'll need to run through several additional steps:

1. "Dockerise" your app - i.e. make it run in a docker container.  This is relatively straight forward, and will simply require you to create a Dockerfile - again, there are examples of these for apps/drivers in each of the directories for each language.
2. Create a databox-manifest.json file.  There are also examples of these in each of the src directories.
3. Build/copy your docker image on the databox.  One simple approach is to register for an account at [docker hub](https://hub.docker.com/) and then push you docker container there.  You can then pull it onto the databox.
4. Name your image so that databox can find it.  By default databox searches for its images at databoxsystems (this is configurable but we'll ignore this for now).  Databox also uses a naming scheme as follows: appname-[architecture]:version. For example if you have an app called myapp, running version 0.5.1 on an 64bit x86-based machine then your image will need to be called:

```
databoxsystems/myapp-amd64:0.5.1
```

to tag your image correctly, simply do the following:

```
docker tag [myimagename] databoxsystems/[myimagename]-amd64:0.5.1
```

5.  Finally, you'll need to upload your manifest file to tell databox about the new app/driver.  Log in to the databox and navigate to My Apps, then click on the "app store" app.  At the top right, use the gear icon to acsess the app store settings form to upload your manifest.  Once uploaded, you can navigate to "App Store" and you should see it ready to install.

The build scripts in the included examples have a number of commands to help with this process, see their README.md's for more details.

## Useful docker commands

Here are some docker commands that help you see what is happening on the databox platform.

| Live | Dev Env | Description |
| --- | --- | --- |
| `docker service logs [app-name] ` | `N/A` | View app/driver console output  |
| `docker service logs arbiter` | `docker logs arbiter` | view the arbiter logs|
| `docker service logs [app-name]-core-store` | `docker logs zest` | view the store logs when running the test env|
| `docker ps` | `docker ps` | check which containers are running|
| `docker service ps` | `N/A` | check which services are running|
| `docker service ps -a [app-name] ` | `N/A` | for debugging service start-up problems if docker docker service logs is empty|
