package logic

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"go.uber.org/fx"
)

func NewDaprClient(lifecycle fx.Lifecycle) (dapr.Client, error) {
	client, err := dapr.NewClient()
	if err != nil {
		return nil, err
	}
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			client.Close()
			return nil
		},
	})
	return client, nil
}
