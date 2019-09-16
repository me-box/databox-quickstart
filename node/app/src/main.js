var https = require("https");
var http = require("http");
var express = require("express");
var bodyParser = require("body-parser");
var databox = require("node-databox");
var WebSocket = require("ws");

const DATABOX_ARBITER_ENDPOINT = process.env.DATABOX_ARBITER_ENDPOINT || 'tcp://127.0.0.1:4444';
const DATABOX_ZMQ_ENDPOINT = process.env.DATABOX_ZMQ_ENDPOINT || "tcp://127.0.0.1:5555";
const DATABOX_TESTING = !(process.env.DATABOX_VERSION);
const PORT = DATABOX_TESTING ? 8090 : process.env.PORT || '8080';

//this will ref the store client which will observe and write to the databox actuator (created in the driver)
let store;
let helloWorldActuatorDataSourceID;

let exportClient = databox.NewExportClient( DATABOX_ARBITER_ENDPOINT, true )
const EXPORT_URL = 'https://postman-echo.com/post'

//server and websocket connection;
let ws, server = null;

const listenToActuator = (emitter) => {

    console.log("started listening to actuator");

    emitter.on('data', (data) => {
        console.log("seen data from the hello world actuator!", data);
        if (ws) {
	    let json = JSON.stringify(data.data)
            ws.send(json);
 	    // Note, export service deprecated and not currently supported
            exportClient.Longpoll( EXPORT_URL, data.data )
	    .then((res) => {
		    console.log('Export ok', res)
		    let poll = function() {
			    exportClient.Longpoll( EXPORT_URL, data.data, res.id )
			    .then((res) => {
				console.log('Export poll ok', res)
				if (res.state !== 'Finished') 
				    setTimeout(poll, 10000)
			    })
			    .catch((err) => {
				    console.log('Export poll error', err)
			    })
		    }
		    if (res.id && res.state !== 'Finished')
			    setTimeout(poll, 1000)
	    })
	    .catch((err) => {
		    console.log('Export error', err)
	    })
	}
    });

    emitter.on('error', (err) => {
        console.warn("error from actuator", err);
    });
}

if (DATABOX_TESTING) {
    store = databox.NewStoreClient(DATABOX_ZMQ_ENDPOINT, DATABOX_ARBITER_ENDPOINT, false);
    helloWorldActuatorDataSourceID = "helloWorldActuator";
    store.TSBlob.Observe(helloWorldActuatorDataSourceID).then((emitter) => {
        listenToActuator(emitter);
    });
} else {
    //listen in on the helloWorld Actuator, which we have asked permissions for in the manifest
    let helloWorldActuator = databox.HypercatToDataSourceMetadata(process.env[`DATASOURCE_helloWorldActuator`]);
    helloWorldActuatorDataSourceID = helloWorldActuator.DataSourceID;
    let helloWorldStore = databox.GetStoreURLFromHypercat(process.env[`DATASOURCE_helloWorldActuator`]);
    store = databox.NewStoreClient(helloWorldStore, DATABOX_ARBITER_ENDPOINT, false)
    store.TSBlob.Observe(helloWorldActuatorDataSourceID, 0)
    .then((emitter) => {
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

    let data = { msg: `${Date.now()}: databox actuation event` };
    return new Promise((resolve, reject) => {
        store.TSBlob.Write(helloWorldActuatorDataSourceID, data).then(() => {
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
    const credentials = databox.GetHttpsCredentials();
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
