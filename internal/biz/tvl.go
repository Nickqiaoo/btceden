package biz

import (
	"btceden/internal/conf"
	"context"
	"strconv"
	"strings"

	v1 "btceden/api/summary/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type TVL struct {
	TVL string
	TPS string
	GAS string
}

type TVLRepo interface {
	TVL(context.Context) ([]*ChainTVL, error)
}

type TVLUsecase struct {
	c       *conf.Data
	tvlRepo TVLRepo
	log     *log.Helper

	tvl map[string]*ChainTVL
}

func NewTVLUsecase(c *conf.Data, tvl TVLRepo, logger log.Logger) *TVLUsecase {
	s := &TVLUsecase{c: c, tvlRepo: tvl, log: log.NewHelper(logger)}
	//s.Load(context.Background())
	return s
}

func (uc *TVLUsecase) TVL(ctx context.Context, chains []string) (res map[string]string, err error) {
	res = make(map[string]string)
	for _, chain := range uc.c.Chains {
		if v, ok := uc.tvl[chain.Name]; ok {
			res[chain.Name] = strconv.FormatFloat(v.Tvl, 'f', -1, 64)
		}
	}
	return
}

func (uc *TVLUsecase) GetTVL() {
	ctx := context.Background()
	res, err := uc.tvlRepo.TVL(ctx)
	if err != nil || res == nil {
		uc.log.WithContext(ctx).Errorf("TVLUsecase GetTVL err: %v", err)
		return
	}
	chains := make(map[string]struct{})
	tvl := make(map[string]*ChainTVL)

	for _, v := range uc.c.Chains {
		chains[v.Name] = struct{}{}
	}
	for _, v := range res {
		if _, ok := chains[strings.ToLower(v.Name)]; ok {
			tvl[strings.ToLower(v.Name)] = v
		}
	}
	uc.tvl = tvl
}

func (uc *TVLUsecase) Load(ctx context.Context) {
	uc.log.WithContext(ctx).Infof("TVLUsecase Load Start")
	uc.GetTVL()
	c := cron.New()
	c.AddFunc("@every 1m", uc.GetTVL)
	//c.AddFunc("0 0 2 * * *", s.LoadWeShineCPIncr)
	//c.AddFunc("@hourly", uc.GetTVL)

	c.Start()
}
