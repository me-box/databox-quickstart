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

	//Are we running inside databox?
	DataboxTestMode := os.Getenv("DATABOX_VERSION") == ""

	//Read in the information on the datasources that databox passed to the app
	var testDataSource libDatabox.DataSourceMetadata
	var storeEndpoint string
	var storeClient *libDatabox.CoreStoreClient
	httpServerPort := "8080"
	if DataboxTestMode {
		libDatabox.Warn("Missing DATASOURCE_TESTDATA assuming we are outside of databox")
		storeEndpoint = testStoreEndpoint
		httpServerPort = "8081" //this is needed to avoid collisions with the driver you can use any free port
		//Fake the datasource information which we would normally get from databox as an env var
		testDataSource = libDatabox.DataSourceMetadata{
			Description:    "A test datasource",        //required
			ContentType:    libDatabox.ContentTypeJSON, //required
			Vendor:         "databox-test",             //required
			DataSourceType: "testData",                 //required
			DataSourceID:   "testdata1",                //required
			StoreType:      libDatabox.StoreTypeTSBlob, //required
			IsActuator:     false,
			IsFunc:         false,
		}
		//turn on debug output for the databox library
		libDatabox.OutputDebug(true)
		//Set up a store client you will need one of these per store
		ac, _ := libDatabox.NewArbiterClient("./", "./", testArbiterEndpoint)
		storeClient = libDatabox.NewCoreStoreClient(ac, "./", storeEndpoint, false)
	} else {
		//This is the standard setup for inside databox
		var err error
		testDataSource, storeEndpoint, err = libDatabox.HypercatToDataSourceMetadata(os.Getenv("DATASOURCE_testdata"))
		libDatabox.ChkErr(err)
		// Set up a store client you will need one of these per store
		// if you asked for more then one data source in your manifest
		// there will be more then one env var provided by databox DATASOURCE_[manifest client id]
		storeClient = libDatabox.NewDefaultCoreStoreClient(storeEndpoint)
	}

	//The endpoints and routing for the app UI
	router := mux.NewRouter()
	router.HandleFunc("/status", statusEndpoint).Methods("GET")
	router.HandleFunc("/ui/getData", getData(testDataSource, storeClient)).Methods("GET")
	router.HandleFunc("/ui/crash", crashApp).Methods("GET")
	router.HandleFunc("/ui/qstest", qstest).Methods("GET")
	router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("./static"))))

	//setup webserver
	setUpWebServer(DataboxTestMode, router, httpServerPort)

	libDatabox.Info("Exiting ....")
}

func qstest(w http.ResponseWriter, r *http.Request) {
	libDatabox.Info(r.URL.Path)
	libDatabox.Info(r.URL.RawPath)
	libDatabox.Info(r.URL.RawQuery)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("active\n"))
}

func crashApp(w http.ResponseWriter, r *http.Request) {
	os.Exit(2)
}

func statusEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("active\n"))
}

func getData(dataSource libDatabox.DataSourceMetadata, store *libDatabox.CoreStoreClient) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//read the latest value from the store
		//the store retunes a json string of the form [{"timestamp":1538464315931,"data":{"data":"44"}}]
		latest, err := store.TSBlobJSON.Latest(dataSource.DataSourceID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":500,"data":"%s"}`, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `%s`, latest)

	}
}

func setUpWebServer(testMode bool, r *mux.Router, port string) {

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
		//Start up a well behaved HTTPS server for displying the UI
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
