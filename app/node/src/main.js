var https = require("https");
var http = require("http");
var express = require("express");
var bodyParser = require("body-parser");
var databox = require("node-databox");
var WebSocket = require("ws");

const DATABOX_ZMQ_ENDPOINT = /*process.env.DATABOX_ZMQ_ENDPOINT ||*/ "tcp://127.0.0.1:5555";
const DATABOX_TESTING = !(process.env.DATABOX_VERSION);
const PORT = DATABOX_TESTING ? 8090 : process.env.port || '8080';

//this will ref the timeseriesblob client which will observe and write to the databoxactuator (created in the driver)
let tsc;

//server and websocket connection;
let ws, server = null;

const listenToActuator = (emitter) => {

    console.log("started listening to actuator");

    emitter.on('data', (data) => {
        console.log("seen data from the hello world actuator!", JSON.parse(data.data));
        if (ws) {
            ws.send(data.data);
        }
    });

    emitter.on('error', (err) => {
        console.warn(err);
    });
}

if (DATABOX_TESTING) {
    tsc = databox.NewTimeSeriesBlobClient(DATABOX_ZMQ_ENDPOINT, false);
    tsc.Observe("helloWorldActuator").then((emitter) => {
        listenToActuator(emitter);
    });
} else {
    let helloWorldActuator;

    //listen in on the helloWorld Actuator, which we have asked permissions for in the manifest
    databox.HypercatToSourceDataMetadata(process.env[`DATASOURCE_helloWorldActuator`]).then((data) => {
        helloWorldActuator = data
        return databox.NewTimeSeriesBlobClient(helloWorldActuator.DataSourceURL, false)
    }).then((store) => {
        tsc = store;
        return store.Observe(helloWorldActuator.DataSourceMetadata.DataSourceID)
    }).then((emitter) => {
        listenToActuator(emitter);
    }).catch((err) => {
        console.warn("Error Observing helloWorldActuator", err);
    });
}


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
    res.render('index', { testing: DATABOX_TESTING });
});

app.get('/ui/actuate', (req, res) => {

    return new Promise((resolve, reject) => {
        tsc.Write("helloWorldActuator", { msg: `${Date.now()}:databox actuation event` }).then(() => {
            console.log("successfully actuated!");
            resolve();
        }).catch((err) => {
            console.log("failed to actuate", err);
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
    server = http.createServer(app).listen(PORT);

} else {
    console.log("[Creating https server]", PORT);
    const credentials = databox.getHttpsCredentials();
    server = https.createServer(credentials, app).listen(PORT);
}

//finally, set up websockets
const wss = new WebSocket.Server({ server, path: "/ui/ws" });

wss.on("connection", (_ws) => {
    ws = _ws;
});

wss.on("error", (err) => {
    console.log("websocket error", err);
    if (ws) {
        ws = null;
    }
})
