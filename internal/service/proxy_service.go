package service

import (
	"btceden/internal/biz"
	"btceden/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type ProxyService struct {
	c     *conf.Data
	log   *log.Helper
	proxy *biz.ProxyUsecase
}

func NewProxyService(c *conf.Data, proxy *biz.ProxyUsecase, logger log.Logger) *ProxyService {
	return &ProxyService{
		c:     c,
		log:   log.NewHelper(logger),
		proxy: proxy,
	}
}

func (s *ProxyService) TVL(ctx context.Context, chainid string) (res map[string]interface{}, err error) {
	return s.proxy.TVL(ctx, chainid)
}
func (s *ProxyService) Activity(ctx context.Context) (res map[string]interface{}, err error) {
	return s.proxy.Activity(ctx)
}
func (s *ProxyService) TVLBreakDown(ctx context.Context, chainid string) (res map[string]interface{}, err error) {
	return s.proxy.TVLBreakDown(ctx, chainid)
}
