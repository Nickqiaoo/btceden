package biz

import (
	"btceden/internal/conf"
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
	"sync"
)

type TPSRepo interface {
	Block(string, string, context.Context) (*Block, error)
}

type Block struct {
	Number    int64
	Timestamp int64
	Tx        int32
}

var NullErr = errors.New("null")

type TPSUsecase struct {
	c       *conf.Data
	tpsRepo TPSRepo
	log     *log.Helper
	cron    *cron.Cron

	number sync.Map
	tps    sync.Map
	queue  map[string]*TPSQueue
}

func NewTPSUsecase(c *conf.Data, tps TPSRepo, logger log.Logger) *TPSUsecase {
	s := &TPSUsecase{c: c, tpsRepo: tps, log: log.NewHelper(logger), cron: cron.New(), queue: make(map[string]*TPSQueue)}
	for _, chain := range c.Chains {
		s.queue[chain.Name] = NewTPSQueue(28800)
	}
	//s.Load(context.Background())
	return s
}

func (uc *TPSUsecase) TPS(ctx context.Context, chains []string) (res map[string]string, err error) {
	res = make(map[string]string)
	for _, chain := range uc.c.Chains {
		if value, ok := uc.tps.Load(chain.Name); ok {
			v := value.(float64)
			res[chain.Name] = fmt.Sprintf("%.1f", v)
		}
	}
	return
}

func (uc *TPSUsecase) GetTPS(chain string) {
	ctx := context.Background()
	var number int64
	if value, ok := uc.number.Load(chain); ok {
		number = value.(int64) + 1
	}
	hexNumber := fmt.Sprintf("0x%x", number)
	if number == 0 {
		hexNumber = "latest"
	}
	res, err := uc.tpsRepo.Block(chain, hexNumber, ctx)
	if err != nil {
		if errors.Is(err, NullErr) {
			return
		}
		uc.log.WithContext(ctx).Errorf("TPSUsecase GetTPS err: %v", err)
		return
	}
	uc.number.Store(chain, res.Number)
	uc.updateTPS(chain, res)
	return
}

func (uc *TPSUsecase) updateTPS(chain string, block *Block) {
	q := uc.queue[chain]
	q.Enqueue(block)
	uc.tps.Store(chain, q.TPS())
}

func (uc *TPSUsecase) InitTPS(chain string) {
	ctx := context.Background()
	res, err := uc.tpsRepo.Block(chain, "latest", ctx)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("TPSUsecase InitTPS err: %v", err)
		return
	}
	uc.number.Store(chain, res.Number)
	uc.queue[chain].setStartTime(res.Timestamp)
	uc.updateTPS(chain, res)
	uc.cron.AddFunc("@every 2s", func() {
		uc.GetTPS(chain)
	})
	return
}

func (uc *TPSUsecase) Load(ctx context.Context) {
	uc.log.WithContext(ctx).Infof("TPSUsecase Load Start")
	for _, v := range uc.c.Chains {
		chainName := v.Name
		go uc.InitTPS(chainName)
	}
	uc.cron.Start()
}

type TPSQueue struct {
	data  []*Block
	head  int
	tail  int
	size  int
	count int

	startTime int64
	endTime   int64
	sum       int64
}

func NewTPSQueue(size int) *TPSQueue {
	return &TPSQueue{
		data: make([]*Block, size),
		size: size,
	}
}

func (cq *TPSQueue) Enqueue(block *Block) {
	if cq.count == cq.size {
		// Queue is full, overwrite the oldest element
		cq.head = (cq.head + 1) % cq.size

		cq.startTime = cq.data[cq.head].Timestamp
		cq.sum -= int64(cq.data[cq.head].Tx)
	} else {
		cq.count++
	}
	cq.data[cq.tail] = block
	cq.tail = (cq.tail + 1) % cq.size

	cq.endTime = block.Timestamp
	cq.sum += int64(block.Tx)
}

func (cq *TPSQueue) setStartTime(start int64) {
	cq.startTime = start
}

func (cq *TPSQueue) TPS() float64 {
	dur := cq.endTime - cq.startTime
	if dur == 0 {
		dur = 3
	}
	tps := float64(cq.sum) / float64(dur)
	return tps
}
