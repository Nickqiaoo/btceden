package service

import (
	"btceden/internal/biz"
	"btceden/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/log"

	pb "btceden/api/summary/v1"
)

type SummaryService struct {
	c *conf.Data
	pb.SummaryHTTPServer
	log *log.Helper
	tvl *biz.TVLUsecase
	tps *biz.TPSUsecase
	gas *biz.GASUsecase
}

func NewSummaryService(c *conf.Data, tvl *biz.TVLUsecase, tps *biz.TPSUsecase, gas *biz.GASUsecase, logger log.Logger) *SummaryService {
	return &SummaryService{
		c:   c,
		log: log.NewHelper(logger),
		tvl: tvl,
		tps: tps,
		gas: gas,
	}
}

func (s *SummaryService) Aggregate(ctx context.Context, req *pb.AggregateRequest) (res *pb.AggregateReply, err error) {
	var (
		tvl map[string]string
		gas map[string]string
		tps map[string]string
	)
	tvl, err = s.tvl.TVL(ctx, []string{req.Chains})
	if err != nil {
		return nil, err
	}
	gas, err = s.gas.GAS(ctx, []string{req.Chains})
	if err != nil {
		return nil, err
	}
	tps, err = s.tps.TPS(ctx, []string{req.Chains})
	if err != nil {
		return nil, err
	}
	res = &pb.AggregateReply{Statistics: make(map[string]*pb.Statistics)}
	for _, chain := range s.c.Chains {
		r := &pb.Statistics{}
		if t, ok := tvl[chain.Name]; ok {
			r.Tvl = t
		}
		if g, ok := gas[chain.Name]; ok {
			r.Gas = g
		}
		if t, ok := tps[chain.Name]; ok {
			r.Tps = t
		}
		res.Statistics[chain.Name] = r
	}
	return
}
