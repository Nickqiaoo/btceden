// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v4.25.3
// source: api/summary/v1/summary.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationSummaryAggregate = "/summary.v1.Summary/Aggregate"

type SummaryHTTPServer interface {
	// Aggregate Sends a greeting
	Aggregate(context.Context, *AggregateRequest) (*AggregateReply, error)
}

func RegisterSummaryHTTPServer(s *http.Server, srv SummaryHTTPServer) {
	r := s.Route("/")
	r.GET("/aggregate", _Summary_Aggregate0_HTTP_Handler(srv))
}

func _Summary_Aggregate0_HTTP_Handler(srv SummaryHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AggregateRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSummaryAggregate)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Aggregate(ctx, req.(*AggregateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AggregateReply)
		return ctx.JSON(200, struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Data    interface{} `json:"data"`
		}{
			Code:    200,
			Message: "OK",
			Data:    reply,
		})
	}
}

type SummaryHTTPClient interface {
	Aggregate(ctx context.Context, req *AggregateRequest, opts ...http.CallOption) (rsp *AggregateReply, err error)
}

type SummaryHTTPClientImpl struct {
	cc *http.Client
}

func NewSummaryHTTPClient(client *http.Client) SummaryHTTPClient {
	return &SummaryHTTPClientImpl{client}
}

func (c *SummaryHTTPClientImpl) Aggregate(ctx context.Context, in *AggregateRequest, opts ...http.CallOption) (*AggregateReply, error) {
	var out AggregateReply
	pattern := "/aggregate"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationSummaryAggregate))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
