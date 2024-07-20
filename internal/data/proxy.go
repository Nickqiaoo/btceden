package data

import (
	"btceden/internal/biz"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-kratos/kratos/v2/log"
)

type proxyRepo struct {
	data   *Data
	log    *log.Helper
	apiUrl string
}

func NewProxyRepo(data *Data, logger log.Logger) biz.ProxyRepo {
	api := os.Getenv("L2BEAT_API_URL")
	if api == "" {
		api = "http://localhost:3000"
	}
	return &proxyRepo{
		data:   data,
		log:    log.NewHelper(logger),
		apiUrl: api,
	}
}

func (r *proxyRepo) Proxy(ctx context.Context, api string) (res map[string]interface{}, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

	if req, err = http.NewRequest("GET", r.apiUrl+api, nil); err != nil {
		r.log.Errorf("proxyRepo http.NewRequest.error(%v)", err)
		return
	}
	if resp, err = r.data.proxy.Do(req); err != nil {
		r.log.Errorf("proxyRepo.Do error(%v)", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.log.Errorf("proxyRepo io.ReadAll() error(%v)", err)
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		r.log.Errorf("proxyRepo json.Unmarshal err (%v)", err)
		return
	}
	if res == nil {
		r.log.Warnf("proxyRepo res(%+v)", res)
		return nil, errors.New(fmt.Sprintf("proxyRepo api(%s) res nil", api))
	}
	return
}
