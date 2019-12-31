package interceptor

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntercept(t *testing.T) {
	ivk := func (ctx context.Context, req interface{}) (interface{},error) {
		fmt.Println("invoker...")
		return nil, nil
	}

	inter1 := func(ctx context.Context, req interface{}, ivk Invoker) (interface{}, error) {
		fmt.Println("interceptor1...")
		return ivk(ctx, req)
	}

	inter2 := func(ctx context.Context, req interface{},  ivk Invoker) (interface{},error) {
		fmt.Println("interceptor2...")
		return ivk(ctx, req)
	}

	ceps := []Interceptor{inter1, inter2}
	_, err := Intercept(context.Background(), nil , ceps , ivk )
	assert.Nil(t, err)
}