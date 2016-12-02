package openshift

import (
	"net/http"

	"github.com/asiainfoLDP/datafoundry_payment/api"
	"github.com/julienschmidt/httprouter"
	"github.com/zonesan/clog"
)

func CreateProject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clog.Info("from", r.RemoteAddr, r.Method, r.URL.RequestURI(), r.Proto)

	org := new(Orgnazition)

	if err := api.ParseRequestBody(r, org); err != nil {
		clog.Error("read request body error.", err)
		api.RespError(w, err)
		return
	}

	oc, err := NewOClient(r)

	if err != nil {
		clog.Error("OpenshiftRestClient", err)
		api.RespError(w, err)
		return
	}

	if proj, err := oc.CreateProject(r, org.Name); err != nil {
		api.RespError(w, err)
	} else {
		api.RespOK(w, proj)
	}

}
