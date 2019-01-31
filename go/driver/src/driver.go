package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	libDatabox "github.com/me-box/lib-go-databox"
)

//default addresses to be used in testing mode
const testArbiterEndpoint = "tcp://127.0.0.1:4444"
const testStoreEndpoint = "tcp://127.0.0.1:5555"

func main() {
	libDatabox.Info("Starting .....")

	//Are we running inside databox?
	DataboxTestMode := os.Getenv("DATABOX_VERSION") == ""

	// Read in the store endpoint provided by databox
	// this is a driver so you will get a core-store
	// and you are responsible for registering datasources
	// and writing in data.
	var DataboxStoreEndpoint string
	var storeClient *libDatabox.CoreStoreClient
	httpServerPort := "8080"
	if DataboxTestMode {
		DataboxStoreEndpoint = testStoreEndpoint
		ac, _ := libDatabox.NewArbiterClient("./", "./", testArbiterEndpoint)
		storeClient = libDatabox.NewCoreStoreClient(ac, "./", DataboxStoreEndpoint, false)
		//turn on debug output for the databox library
		libDatabox.OutputDebug(true)
	} else {
		DataboxStoreEndpoint = os.Getenv("DATABOX_ZMQ_ENDPOINT")
		storeClient = libDatabox.NewDefaultCoreStoreClient(DataboxStoreEndpoint)
	}

	// start a go routine to do some long running work.
	// You can have may of these structure you program
	// as you see fit.
	go doDriverWork(DataboxTestMode, storeClient)

	//The endpoints and routing for the UI
	router := mux.NewRouter()
	router.HandleFunc("/status", statusEndpoint).Methods("GET")
	router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("./static"))))
	setUpWebServer(DataboxTestMode, router, httpServerPort)

	libDatabox.Info("Exiting ....")

}

func doDriverWork(testMode bool, storeClient *libDatabox.CoreStoreClient) {
	libDatabox.Info("starting doDriverWork")

	//register our datasources
	//we only need to do this once at start up
	testDatasource := libDatabox.DataSourceMetadata{
		Description:    "A test datasource",        //required
		ContentType:    libDatabox.ContentTypeJSON, //required
		Vendor:         "databox-test",             //required
		DataSourceType: "databox-test:testdata",    //required
		DataSourceID:   "testdata1",                //required
		StoreType:      libDatabox.StoreTypeTSBlob, //required
		IsActuator:     false,
		IsFunc:         false,
	}
	err := storeClient.RegisterDatasource(testDatasource)
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
		err := storeClient.TSBlobJSON.Write("testdata1", []byte(jsonData))
		if err != nil {
			libDatabox.Err("Error Write Datasource " + err.Error())
		}
		libDatabox.Info("Data written to store testing: " + jsonData)
		time.Sleep(time.Second * 1)
	}
}

func statusEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("active\n"))
}

func setUpWebServer(testMode bool, r *mux.Router, port string) {

	//Start up a well behaved HTTP/S server for displying the UI

	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      r,
	}

	if testMode {
		//set up an http server for testing
		libDatabox.Info("Waiting for http requests on port http://127.0.0.1" + srv.Addr + "/ui ....")
		log.Fatal(srv.ListenAndServe())
	} else {
		//configure tls
		tlsConfig := &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
			},
		}

		srv.TLSConfig = tlsConfig

		libDatabox.Info("Waiting for https requests on port " + srv.Addr + " ....")
		log.Fatal(srv.ListenAndServeTLS(libDatabox.GetHttpsCredentials(), libDatabox.GetHttpsCredentials()))
	}
}
