package data

import (
	"btceden/internal/biz"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"math/big"
	"net/http"
)

type tpsRepo struct {
	data *Data
	rpc  map[string]string
	log  *log.Helper
}

func NewTPSRepo(data *Data, logger log.Logger) biz.TPSRepo {
	r := &tpsRepo{
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

type EthBlockResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  *struct {
		Number       string   `json:"number"`
		Timestamp    string   `json:"timestamp"`
		Transactions []string `json:"transactions"`
	} `json:"result"`
}

func (r *tpsRepo) Block(chain, param string, ctx context.Context) (block *biz.Block, err error) {
	httpClient := r.data.chains[chain]
	var (
		req  *http.Request
		resp *http.Response
	)
	reqData := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{param, false},
		"id":      1,
	}
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		r.log.Errorf("chain(%s) TPS Failed to marshal request data: %v", chain, err)
		return nil, err
	}
	if req, err = http.NewRequest("POST", r.rpc[chain], bytes.NewBuffer(reqBody)); err != nil {
		r.log.Errorf("chain(%s) TPS http.NewRequest.error(%v)", chain, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if resp, err = httpClient.Do(req); err != nil {
		r.log.Errorf("chain(%s) TPS.Do error(%v)", chain, err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.log.Errorf("chain(%s) TPS io.ReadAll() error(%v)", chain, err)
		return nil, err
	}
	var respData EthBlockResponse
	err = json.Unmarshal(body, &respData)
	if err != nil {
		r.log.Errorf("chain(%s) TPS json.Unmarshal err (%+v)", chain, err)
		return nil, err
	}

	if respData.Result == nil {
		return nil, biz.NullErr
	}

	block = &biz.Block{}
	bigInt := new(big.Int)
	_, success := bigInt.SetString(respData.Result.Timestamp, 0)
	if !success {
		r.log.Errorf("chain(%s) TPS Failed to convert Timestamp to big.Int", chain)
		return nil, errors.New("failed to convert Timestamp")
	}
	block.Timestamp = bigInt.Int64()

	_, success = bigInt.SetString(respData.Result.Number, 0)
	if !success {
		r.log.Errorf("chain(%s) TPS Failed to convert Number to big.Int", chain)
		return nil, errors.New("failed to convert Number")
	}
	block.Number = bigInt.Int64()
	block.Tx = int32(len(respData.Result.Transactions))
	return
}
