package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"net/http"

	"github.com/gorilla/mux"

	libDatabox "github.com/me-box/lib-go-databox"
)

//default addresses to be used in testing mode
const testArbiterEndpoint = "tcp://127.0.0.1:4444"
const testStoreEndpoint = "tcp://127.0.0.1:5555"

func main() {
	libDatabox.Info("Starting ....")

	//turn on debug output for the databox library
	libDatabox.OutputDebug(true)

	//Read in the store endpoint provided by databox
	DataboxZmqEndpoint := os.Getenv("DATABOX_STORE_ENDPOINT")
	DataboxTestMode := false
	if DataboxZmqEndpoint == "" {
		//if this is not available we are outside databox or something ver bad has happened
		libDatabox.Warn("DATABOX_ZMQ_ENDPOINT is not set running in test mode")
		DataboxZmqEndpoint = testStoreEndpoint
		DataboxTestMode = true
	}

	go doDriverWork(DataboxTestMode, DataboxZmqEndpoint)
	setUpWebServer(DataboxTestMode)

}

func doDriverWork(testMode bool, storeEndpoint string) {
	libDatabox.Info("starting doDriverWork")
	//connect to the store
	var coreStoreClient *libDatabox.CoreStoreClient
	if testMode {
		ac, _ := libDatabox.NewArbiterClient("./", "./", "tcp://127.0.0.1:4444")
		coreStoreClient = libDatabox.NewCoreStoreClient(ac, "./", storeEndpoint, false)
	} else {
		coreStoreClient = libDatabox.NewDefaultCoreStoreClient(storeEndpoint)
	}
	libDatabox.Info("Connected to store")

	//register our datasources
	//we only need to do this once at start up
	testDatasource := libDatabox.DataSourceMetadata{
		Description:    "A test datasource",        //required
		ContentType:    libDatabox.ContentTypeJSON, //required
		Vendor:         "databox-test",             //required
		DataSourceType: "testData",                 //required
		DataSourceID:   "testdata1",                //required
		StoreType:      libDatabox.StoreTypeTSBlob, //required
		IsActuator:     false,
		IsFunc:         false,
	}
	err := coreStoreClient.RegisterDatasource(testDatasource)
	if err != nil {
		libDatabox.Err("Error Registering Datasource " + err.Error())
		return
	}
	libDatabox.Info("Registered Datasource")

	//do some work forever and write data to the store
	writeCount := int64(0)
	for {
		writeCount++
		jsonData := fmt.Sprintf(`{"data":"%d"}`, writeCount)
		err := coreStoreClient.TSBlobJSON.Write("testdata1", []byte(jsonData))
		if err != nil {
			libDatabox.Err("Error Write Datasource " + err.Error())
		}
		libDatabox.Info("Data written to store: " + jsonData)
		time.Sleep(time.Second * 1)
	}
}

func setUpWebServer(testMode bool) {
	//setup webserver routes
	router := mux.NewRouter()

	//The endpoints and routing for the app
	router.HandleFunc("/status", statusEndpoint).Methods("GET")
	router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("./static"))))

	if testMode {
		//set up an http server fr testing
		srv := &http.Server{
			Addr:         ":8080",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  30 * time.Second,
			Handler:      router,
		}
		libDatabox.Info("Waiting for http requests on port http://127.0.0.1:8080/ui ....")
		log.Fatal(srv.ListenAndServe())
	} else {
		//Start up a well behaved HTTPS server for displying the UI
		tlsConfig := &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
			},
		}

		srv := &http.Server{
			Addr:         ":8080",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  30 * time.Second,
			TLSConfig:    tlsConfig,
			Handler:      router,
		}
		libDatabox.Info("Waiting for https requests on port 8080 ....")
		log.Fatal(srv.ListenAndServeTLS(libDatabox.GetHttpsCredentials(), libDatabox.GetHttpsCredentials()))
	}
	libDatabox.Info("Exiting ....")
}

func statusEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("active\n"))
}

func getData(w http.ResponseWriter, r *http.Request) {
	//todo
}
