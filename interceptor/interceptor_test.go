package interceptor

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntercept(t *testing.T) {
	ivk := func (ctx context.Context, req interface{}, rsp interface{}) error {
		fmt.Println("invoker...")
		return nil
	}

	inter1 := func(ctx context.Context, req interface{}, rsp interface{}, ivk Invoker) error {
		fmt.Println("interceptor1...")
		return ivk(ctx, req, rsp)
	}

	inter2 := func(ctx context.Context, req interface{}, rsp interface{}, ivk Invoker) error {
		fmt.Println("interceptor2...")
		return ivk(ctx, req, rsp)
	}

	ceps := []Interceptor{inter1, inter2}
	assert.Nil(t,Intercept(context.Background(), nil , nil , ceps , ivk ))
}