## Writing an app in node
This app will write a message to the helloworld actuator, created in the node driver code ([home]/databox-data-tracker/driver/node/src). To get started, first ensure that your are running the test environment (i.e. you have called [home]/databox-quickstart/testenv/start.sh) and that you are running the driver code.  Then do:

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