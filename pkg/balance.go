package pkg

import (
	"encoding/json"
	"net/http"
	"net/url"
	// "reflect"

	"github.com/zonesan/clog"
)

//type BalanceAgent service

type BalanceAgent struct {
	*Agent
	BaseURL *url.URL
}

// type Balance struct {
// 	apiBalance
// }

type Balance apiBalance

func (agent *BalanceAgent) Get(r *http.Request) (*Balance, error) {
	// balance := &Balance{
	// 	Balance: 50000.89,
	// 	Status:  "active",
	// }
	// return balance

	urlStr := "/charge/v1/balance"

	if r.URL.RawQuery != "" {
		urlStr += "?" + r.URL.RawQuery
	}

	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := agent.BaseURL.ResolveReference(rel)

	token, err := getToken(r)
	if err != nil {
		clog.Error(err)
		return nil, err
	}

	req, err := agent.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)

	response := new(RemoteResponse)

	if err := agent.Do(req, response); err != nil {
		return nil, err
	}

	balance := new(Balance)

	if err := json.Unmarshal([]byte(response.Data), balance); err != nil {
		clog.Error(err)
		return nil, err
	}

	return balance, nil

}
