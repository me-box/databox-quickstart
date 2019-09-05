var https = require("https");
var http = require("http");
var express = require("express");
var bodyParser = require("body-parser");
var databox = require("node-databox");

const DATABOX_ARBITER_ENDPOINT = process.env.DATABOX_ARBITER_ENDPOINT || 'tcp://127.0.0.1:4444';
const DATABOX_ZMQ_ENDPOINT = process.env.DATABOX_ZMQ_ENDPOINT || "tcp://127.0.0.1:5555";
const DATABOX_TESTING = !(process.env.DATABOX_VERSION);
const PORT = process.env.port || '8080';

const store = databox.NewStoreClient(DATABOX_ZMQ_ENDPOINT, DATABOX_ARBITER_ENDPOINT);

//get the default store metadata
const metaData = databox.NewDataSourceMetadata();

//create store schema for saving key/value config data
const helloWorldConfig = {
    ...databox.NewDataSourceMetadata(),
    Description: 'hello world config',
    ContentType: 'application/json',
    Vendor: 'Databox Inc.',
    DataSourceType: 'helloWorldConfig',
    DataSourceID: 'helloWorldConfig',
    StoreType: 'kv',
}

//create store schema for an actuator (i.e a store that can be written to by an app)
const helloWorldActuator = {
    ...metaData,
    Description: 'hello world actuator',
    ContentType: 'application/json',
    Vendor: 'Databox Inc.',
    DataSourceType: 'helloWorldActuator',
    DataSourceID: 'helloWorldActuator',
    StoreType: 'ts/blob',
    IsActuator: true,
}

///now create our stores using our clients.
store.RegisterDatasource(helloWorldConfig).then(() => {
    console.log("registered helloWorldConfig");
    //now register the actuator
    return store.RegisterDatasource(helloWorldActuator)
}).catch((err) => { console.log("error registering helloWorld config datasource", err) }).then(() => {
    console.log("registered helloWorldActuator, observing", helloWorldActuator.DataSourceID);
    store.TSBlob.Observe(helloWorldActuator.DataSourceID, 0)
        .catch((err) => {
            console.log("[Actuation observing error]", err);
        })
        .then((eventEmitter) => {
            if (eventEmitter) {
                eventEmitter.on('data', (data) => {
                    console.log("[Actuation] data received ", data);
                });
            }
        })
        .catch((err) => {
            console.log("[Actuation error]", err);
        });
});

//set up webserver to serve driver endpoints
const app = express();
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());
app.set('views', './views');
app.set('view engine', 'ejs');

app.get("/", function (req, res) {
    res.redirect("/ui");
});

app.get("/ui", function (req, res) {
    store.KV.Read(helloWorldConfig.DataSourceID, "config").then((result) => {
        console.log("result:", helloWorldConfig.DataSourceID, result);
        res.render('index', { config: result.value });
    }).catch((err) => {
        console.log("get config error", err);
        res.send({ success: false, err });
    });
});

app.post('/ui/setConfig', (req, res) => {

    const config = req.body.config;

    return new Promise((resolve, reject) => {
        store.KV.Write(helloWorldConfig.DataSourceID, "config", { key: helloWorldConfig.DataSourceID, value: config }).then(() => {
            console.log("successfully written!", config);
            resolve();
        }).catch((err) => {
            console.log("failed to write", err);
            reject(err);
        });
    }).then(() => {
        res.send({ success: true });
    });
});

app.get("/status", function (req, res) {
    res.send("active");
});

//when testing, we run as http, (to prevent the need for self-signed certs etc);
if (DATABOX_TESTING) {
    console.log("[Creating TEST http server]", PORT);
    http.createServer(app).listen(PORT);
} else {
    console.log("[Creating https server]", PORT);
    const credentials = databox.GetHttpsCredentials();
    https.createServer(credentials, app).listen(PORT);
}
