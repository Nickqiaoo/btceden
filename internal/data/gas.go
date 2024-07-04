package data

import (
	"btceden/internal/biz"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"math/big"
	"net/http"
)

type gasRepo struct {
	data *Data
	rpc  map[string]string
	log  *log.Helper
}

type Request struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type Response struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

func NewGASRepo(data *Data, logger log.Logger) biz.GASRepo {
	r := &gasRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
	rpc := make(map[string]string)
	for _, v := range data.c.Chains {
		rpc[v.Name] = v.Rpc
	}
	r.rpc = rpc
	return r
}

func (r *gasRepo) GAS(chain string, ctx context.Context) (gas string, err error) {
	httpClient := r.data.chains[chain]
	var (
		req  *http.Request
		resp *http.Response
	)
	reqData := Request{
		JSONRPC: "2.0",
		Method:  "eth_gasPrice",
		Params:  []interface{}{},
		ID:      1,
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		r.log.Errorf("chain(%s) GAS Failed to marshal request data: %v", chain, err)
	}
	if req, err = http.NewRequest("POST", r.rpc[chain], bytes.NewBuffer(reqBody)); err != nil {
		r.log.Errorf("chain(%s) GAS http.NewRequest.error(%v)", chain, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	if resp, err = httpClient.Do(req); err != nil {
		r.log.Errorf("chain(%s) GAS.Do error(%v)", chain, err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.log.Errorf("chain(%s) GAS io.ReadAll() error(%v)", chain, err)
		return
	}
	var respData Response
	err = json.Unmarshal(body, &respData)
	if err != nil {
		r.log.Errorf("chain(%s) GAS json.Unmarshal err (%v)", chain, err)
		return
	}

	bigInt := new(big.Int)
	_, success := bigInt.SetString(respData.Result, 0)
	if !success {
		r.log.Errorf("chain(%s) GAS Failed to convert hex to big.Int", chain)
		return
	}

	gas = bigInt.String()
	return
}
