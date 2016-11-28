package pkg

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/zonesan/clog"
)

type RechargeAgent struct {
	*Agent
	BaseURL *url.URL
}

type Recharge struct {
	Amount  float64 `json:"amount,omitempty"`
	Project string  `json:"namespace,omitempty"`
}

type HongPay apiRechargePayload

func (agent *RechargeAgent) Create(r *http.Request, recharge *Recharge) (*HongPay, error) {

	urlStr := "/charge/v1/recharge"

	hongpay := new(HongPay)

	if err := doRequest(agent, r, "POST", urlStr, recharge, hongpay); err != nil {
		clog.Error(err)

		return nil, err
	}
	clog.Debug(hongpay.Payloads)
	return hongpay, nil

}

func (agent *RechargeAgent) Notification(r *http.Request) ([]byte, error) {
	urlStr := "/charge/v1/aipaycallback"

	if r.URL.RawQuery != "" {
		urlStr += "?" + r.URL.RawQuery
	}

	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := agent.BaseURL.ResolveReference(rel)

	reqbody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		return nil, err
	}

	clog.Debug("Request Body:", string(reqbody))

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(reqbody))
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	resp, err := agent.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	data, err := ioutil.ReadAll(resp.Body)

	clog.Debugf("%s", data)

	return data, err

}

func (agent *RechargeAgent) Url() *url.URL {
	u := new(url.URL)
	u, _ = url.Parse(agent.BaseURL.String())
	return u
}

func (agent *RechargeAgent) Instance() *Agent {
	return agent.Agent
}
