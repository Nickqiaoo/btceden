package biz

import (
	"btceden/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/mohae/deepcopy"
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
	"200901":   "bitlayer",
	"4200":     "merlin",
	"30":       "rootstock",
	"11501":    "bevm",
	"60808":    "bob",
	"223":      "bsquared",
	"10100001": "stacks",
	"10100002": "liquid",
	"3109":     "satoshivm",
	"57":       "syscoin",
	"10100003": "ckbtc",
	"10100004": "libre",
	"1456":     "zkbase",
	"10100005": "ark",
	"1116":     "core",
	"6001":     "bouncebit",
	"2649":     "ailayer",
	"10100006": "lightning",
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
	if value, ok := uc.data.Load("tvl-project"); ok {
		res = value.(map[string]interface{})
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

func (uc *ProxyUsecase) Layer2sTVL(ctx context.Context) (res map[string]interface{}) {
	if value, ok := uc.data.Load("tvl-layer2s"); ok {
		res = value.(map[string]interface{})
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
	timestamp := res["dataTimestamp"]
	if projects, exists := res["breakdowns"].(map[string]interface{}); exists {
		if p, exist := projects[project].(map[string]interface{}); exist {
			p["dataTimestamp"] = timestamp
			return p, nil
		}
	}
	return
}

func (uc *ProxyUsecase) Activity(ctx context.Context, chainid string) (res map[string]interface{}, err error) {
	if value, ok := uc.data.Load("activity-project"); ok {
		res = value.(map[string]interface{})
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

func (uc *ProxyUsecase) Layer2sActivity(ctx context.Context) (res map[string]interface{}) {
	if value, ok := uc.data.Load("activity-layer2s"); ok {
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
	copied := deepcopy.Copy(res).(map[string]interface{})
	delete(copied, "layers2s")
	uc.data.Store("tvl-project", copied)

	if projects, exists := res["projects"].(map[string]interface{}); exists {
		for _, chain := range projects {
			if c, isc := chain.(map[string]interface{}); isc {
				delete(c, "charts")
			}
		}

	}
	uc.data.Store("tvl-layer2s", res)
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
	copied := deepcopy.Copy(res).(map[string]interface{})
	delete(copied, "combined")
	uc.data.Store("activity-project", copied)

	if projects, exists := res["projects"].(map[string]interface{}); exists {
		for _, chain := range projects {
			if c, isc := chain.(map[string]interface{}); isc {
				delete(c, "daily")
			}
		}

	}
	uc.data.Store("activity-layer2s", res)
	return
}

func (uc *ProxyUsecase) Load(ctx context.Context) {
	uc.log.WithContext(ctx).Infof("ProxyUsecase Load Start")
	c := cron.New()
	uc.getTVL()
	uc.getTVLBreakDown()
	uc.getActivity()

	c.AddFunc("10 * * * *", uc.getTVL)
	c.AddFunc("10 * * * *", uc.getTVLBreakDown)
	c.AddFunc("@every 1h", uc.getActivity)

	c.Start()
}
