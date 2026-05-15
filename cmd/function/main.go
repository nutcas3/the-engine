package main

import (
	"context"

	"the-engine/internal/function"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
)

type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer
	handler *function.Handler
}

func NewFunction() *Function {
	return &Function{
		handler: &function.Handler{},
	}
}

func (f *Function) RunFunction(ctx context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	return f.handler.RunFunction(ctx, req)
}
