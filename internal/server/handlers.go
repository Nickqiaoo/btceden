package server

import (
	"btceden/internal/service"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var proxyService *service.ProxyService

func tvl(ctx http.Context) error {
	res, err := proxyService.TVL(ctx)
	if err != nil {
		return err
	}

	return ctx.JSON(200, struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Code:    200,
		Message: "OK",
		Data:    res,
	})
}

func breakdown(ctx http.Context) error {
	res, err := proxyService.TVLBreakDown(ctx)
	if err != nil {
		return err
	}
	return ctx.JSON(200, struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Code:    200,
		Message: "OK",
		Data:    res,
	})
}

func activity(ctx http.Context) error {
	res, err := proxyService.Activity(ctx)
	if err != nil {
		return err
	}
	return ctx.JSON(200, struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Code:    200,
		Message: "OK",
		Data:    res,
	})
}
