package server

import "github.com/go-kratos/kratos/v2/transport/http"

func initRoute(svr *http.Server) {
	r := svr.Route("/")
	r2 := r.Group("/api")
	{
		r2.GET("/tvl", tvl)
		r2.GET("/tvl/breakdown", breakdown)
		r2.GET("activity", activity)
	}
}
