package interceptor

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntercept(t *testing.T) {
	ivk := func (ctx context.Context, req, rsp interface{}) error {
		fmt.Println("invoker...")
		return nil
	}

	inter1 := func(ctx context.Context, req, rsp interface{}, ivk Invoker) error {
		fmt.Println("interceptor1...")
		return ivk(ctx, req,rsp)
	}

	inter2 := func(ctx context.Context, req,rsp interface{},  ivk Invoker) error {
		fmt.Println("interceptor2...")
		return ivk(ctx, req,rsp)
	}
	ceps := []ClientInterceptor{inter1, inter2}

	err := ClientIntercept(context.Background(), nil ,nil, ceps , ivk)
	assert.Nil(t, err)
}