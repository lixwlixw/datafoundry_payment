package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	payment "github.com/asiainfoLDP/datafoundry_payment/pkg"
	"github.com/zonesan/clog"
)

func ParseRequestBody(r *http.Request, v interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	clog.Debug("Request Body:", string(b))
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

func RespError(w http.ResponseWriter, err error) {
	resp := genRespJson(err)

	if body, err := json.MarshalIndent(resp, "", "  "); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.status)
		w.Write(body)
	}

}

func RespOK(w http.ResponseWriter, data interface{}) {
	if data == nil {
		data = genRespJson(nil)
	}

	if body, err := json.MarshalIndent(data, "", "  "); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func genRespJson(err error) *APIResponse {
	resp := new(APIResponse)

	if err == nil {
		resp.Code = payment.ErrCodeOK
		resp.status = http.StatusOK
	} else {
		if e, ok := err.(*payment.ErrorMessage); ok {
			resp.Code = e.Code
			resp.status = trickCode2Status(resp.Code) //http.StatusBadRequest
			resp.Message = payment.ErrText(resp.Code)
		} else if e, ok := err.(*payment.ErrorResponse); ok {
			//TODO
			resp.Code = e.Code
			resp.Message = e.Message
			resp.status = e.Response.StatusCode
		} else {
			resp.Code = payment.ErrCodeBadRequest
			resp.Message = err.Error()
			resp.status = trickCode2Status(resp.Code) //http.StatusBadRequest
		}
	}

	resp.Reason = http.StatusText(resp.status)

	return resp
}

func trickCode2Status(errCode int) int {
	var statusCode int
	if errCode < 10000 {
		statusCode = errCode % 1000
	} else {
		statusCode = trickCode2Status(errCode / 10)
	}

	return statusCode
}
