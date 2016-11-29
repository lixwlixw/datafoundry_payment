package main

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/pkg/openshift"
	"github.com/zonesan/clog"
)

func main() {

	router := createRouter()

	//clog.SetLogLevel(clog.LOG_LEVEL_DEBUG)
	openshift.Hello()
	clog.Info("listening on port 8080...")
	clog.Fatal(http.ListenAndServe(":8080", router))
}
