package pkg

import (
	//"encoding/json"
	"net/http"
	"net/url"

	"github.com/zonesan/clog"
)

type CheckoutAgent struct {
	*Agent
	BaseURL *url.URL
}

type Checkout struct {
	PlanId  string                 `json:"plan_id"`
	Project string                 `json:"namespace,omitempty"`
	Region  string                 `json:"region"` //need it?
	OrderID string                 `json:"orderid,omitempty"`
	DryTry  int                    `json:"drytry,omitempty"`
	Params  map[string]interface{} `json:"parameters,omitempty"`
}

type PurchasedOrder apiPurchaseOrder

func (agent *CheckoutAgent) Get() *Balance {

	balance := &Balance{
		Balance: 3000.01,
		Status:  "active",
	}
	return balance
}

func (agent *CheckoutAgent) ListOrders(r *http.Request) (*[]PurchasedOrder, error) {
	urlStr := "/usageapi/v1/orders"

	orders := new([]PurchasedOrder)

	if err := doRequestList(agent, r, "GET", urlStr, nil, orders); err != nil {
		clog.Error(err)
		return nil, err
	}

	clog.Infof("%v order(s) listed.", len(*orders))

	return orders, nil

	// if r.URL.RawQuery != "" {
	// 	urlStr += "?" + r.URL.RawQuery
	// }

	// rel, err := url.Parse(urlStr)
	// if err != nil {
	// 	return nil, err
	// }

	// u := agent.BaseURL.ResolveReference(rel)

	// token, err := getToken(r)
	// if err != nil {
	// 	clog.Error(err)
	// 	return nil, err
	// }

	// req, err := agent.NewRequest("GET", u.String(), nil)
	// if err != nil {
	// 	return nil, err
	// }

	// req.Header.Set("Authorization", token)

	// response := new(RemoteListResponse)
	// if err := agent.Do(req, response); err != nil {
	// 	clog.Error(err)
	// 	return nil, err
	// }

	// orders := []PurchasedOrder{}

	// if err := json.Unmarshal([]byte(response.Data), &orders); err != nil {
	// 	clog.Error(err)
	// 	return nil, err
	// } else {
	// 	clog.Infof("%v order(s) listed.", len(orders))
	// 	return &orders, nil
	// }

}

func (agent *CheckoutAgent) Create(r *http.Request, checkout *Checkout) (*PurchasedOrder, error) {
	urlStr := "/usageapi/v1/orders"

	order := new(PurchasedOrder)
	if err := doRequest(agent, r, "POST", urlStr, checkout, order); err != nil {
		clog.Error(err)
		return nil, err
	}

	return order, nil
}

func (agent *CheckoutAgent) Url() *url.URL {
	u := new(url.URL)
	u, _ = url.Parse(agent.BaseURL.String())
	return u
}

func (agent *CheckoutAgent) Instance() *Agent {
	return agent.Agent
}
