package biz

import (
	"btceden/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
	"sync"
)

type GASRepo interface {
	GAS(string, context.Context) (string, error)
}

type GASUsecase struct {
	c       *conf.Data
	gasRepo GASRepo
	log     *log.Helper

	gas sync.Map
}

func NewGASUsecase(c *conf.Data, gas GASRepo, logger log.Logger) *GASUsecase {
	s := &GASUsecase{c: c, gasRepo: gas, log: log.NewHelper(logger)}
	//s.Load(context.Background())
	return s
}

func (uc *GASUsecase) GAS(ctx context.Context, chains []string) (res map[string]string, err error) {
	res = make(map[string]string)
	for _, chain := range uc.c.Chains {
		if value, ok := uc.gas.Load(chain.Name); ok {
			res[chain.Name] = value.(string)
		}
	}
	return
}

func (uc *GASUsecase) GetGAS(chain string) {
	ctx := context.Background()
	res, err := uc.gasRepo.GAS(chain, ctx)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("GASUsecase GetGAS err: %v", err)
		return
	}
	uc.gas.Store(chain, res)
	return
}

func (uc *GASUsecase) Load(ctx context.Context) {
	uc.log.WithContext(ctx).Infof("GASUsecase Load Start")
	c := cron.New()
	for _, v := range uc.c.Chains {
		chainName := v.Name
		c.AddFunc("@every 1m", func() {
			uc.GetGAS(chainName)
		})
	}
	//c.AddFunc("0 0 2 * * *", s.LoadWeShineCPIncr)
	//c.AddFunc("@hourly", uc.GetTVL)

	c.Start()
}
