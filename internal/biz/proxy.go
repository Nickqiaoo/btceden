package biz

import (
	"btceden/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
	"sync"
)

type ProxyUsecase struct {
	c         *conf.Data
	proxyRepo ProxyRepo
	log       *log.Helper

	data sync.Map
}

var chainIdMap = map[string]string{
	"200901": "bitlayer",
	"4200":   "Merlin",
	"30":     "Rootstock",
	"11501":  "BEVM",
	"60808":  "BOB",
}

type ProxyRepo interface {
	Proxy(context.Context, string) (map[string]interface{}, error)
}

func NewProxyUsecase(c *conf.Data, proxyRepo ProxyRepo, logger log.Logger) *ProxyUsecase {
	s := &ProxyUsecase{c: c, proxyRepo: proxyRepo, log: log.NewHelper(logger)}
	s.Load(context.Background())
	return s
}

func (uc *ProxyUsecase) TVL(ctx context.Context, chainid string) (res map[string]interface{}, err error) {
	if value, ok := uc.data.Load("tvl"); ok {
		res = value.(map[string]interface{})
	}
	if chainid == "" {
		return
	}
	var (
		project string
		exist   bool
	)
	if project, exist = chainIdMap[chainid]; !exist {
		return
	}
	if projects, exists := res["projects"].(map[string]interface{}); exists {
		if p, exist := projects[project].(map[string]interface{}); exist {
			return p, nil
		}
	}
	return
}

func (uc *ProxyUsecase) TVLBreakDown(ctx context.Context, chainid string) (res map[string]interface{}, err error) {
	if value, ok := uc.data.Load("tvl-breakdown"); ok {
		res = value.(map[string]interface{})
	}
	var (
		project string
		exist   bool
	)
	if project, exist = chainIdMap[chainid]; !exist {
		return
	}
	if projects, exists := res["breakdowns"].(map[string]interface{}); exists {
		if p, exist := projects[project].(map[string]interface{}); exist {
			return p, nil
		}
	}
	return
}

func (uc *ProxyUsecase) Activity(ctx context.Context) (res map[string]interface{}, err error) {
	if value, ok := uc.data.Load("activity"); ok {
		res = value.(map[string]interface{})
	}
	return
}

func (uc *ProxyUsecase) getTVL() {
	ctx := context.Background()
	res, err := uc.proxyRepo.Proxy(ctx, "/api/tvl")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("ProxyUsecase getTVL err: %v", err)
		return
	}
	uc.data.Store("tvl", res)
	return
}

func (uc *ProxyUsecase) getTVLBreakDown() {
	ctx := context.Background()
	res, err := uc.proxyRepo.Proxy(ctx, "/api/project-assets-breakdown")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("ProxyUsecase getTVLBreakDown err: %v", err)
		return
	}
	uc.data.Store("tvl-breakdown", res)
	return
}

func (uc *ProxyUsecase) getActivity() {
	ctx := context.Background()
	res, err := uc.proxyRepo.Proxy(ctx, "/api/activity")
	if err != nil {
		uc.log.WithContext(ctx).Errorf("ProxyUsecase getActivity err: %v", err)
		return
	}
	uc.data.Store("activity", res)
	return
}

func (uc *ProxyUsecase) Load(ctx context.Context) {
	uc.log.WithContext(ctx).Infof("ProxyUsecase Load Start")
	c := cron.New()
	uc.getTVL()
	uc.getTVLBreakDown()
	uc.getActivity()

	c.AddFunc("@every 1h", uc.getTVL)
	c.AddFunc("@every 1h", uc.getTVLBreakDown)
	c.AddFunc("@every 1h", uc.getActivity)

	c.Start()
}
