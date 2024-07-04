package data

import (
	"btceden/internal/biz"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
)

type tvlRepo struct {
	data *Data
	log  *log.Helper
}

func NewTVLRepo(data *Data, logger log.Logger) biz.TVLRepo {
	return &tvlRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *tvlRepo) TVL(ctx context.Context) (res []*biz.ChainTVL, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

	if req, err = http.NewRequest("GET", r.data.c.DefiLlamaApi, nil); err != nil {
		r.log.Errorf("tvlRepo http.NewRequest.error(%v)", err)
		return
	}
	if resp, err = r.data.defillama.Do(req); err != nil {
		r.log.Errorf("tvlRepo.Do error(%v)", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.log.Errorf("tvlRepo io.ReadAll() error(%v)", err)
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		r.log.Errorf("tvlRepo json.Unmarshal err (%v)", err)
		return
	}
	if res == nil {
		r.log.Warnf("tvlRepo res(%+v)", res)
		return
	}
	return
}
