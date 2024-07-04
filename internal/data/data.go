package data

import (
	"btceden/internal/conf"
	"context"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewTVLRepo, NewTPSRepo, NewGASRepo, NewProxyRepo)

// Data .
type Data struct {
	c         *conf.Data
	defillama *http.Client
	chains    map[string]*http.Client
	proxy     *http.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	var err error
	client, err := http.NewClient(context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithTimeout(time.Second*5))
	if err != nil {
		return nil, cleanup, err
	}
	client1, err := http.NewClient(context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithTimeout(time.Second*5))
	if err != nil {
		return nil, cleanup, err
	}
	chainRpc := make(map[string]*http.Client)
	for _, v := range c.Chains {
		chainRpc[v.Name], err = http.NewClient(context.Background(),
			http.WithMiddleware(
				recovery.Recovery(),
			),
			http.WithTimeout(time.Second*20))
		if err != nil {
			return nil, cleanup, err
		}
	}
	return &Data{c: c, defillama: client, chains: chainRpc, proxy: client1}, cleanup, nil
}
