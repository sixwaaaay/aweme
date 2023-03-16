package main

import (
	"github.com/PlanVX/aweme/internal/config"
	"github.com/PlanVX/aweme/internal/dal"
	"github.com/PlanVX/aweme/internal/dal/query"
	"github.com/PlanVX/aweme/internal/otel"
	"github.com/PlanVX/aweme/internal/sub"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	var dataLayer = []any{
		query.NewGormDB,
		query.NewRedisUniversalClient,
		fx.Annotate(query.NewCommentCommand, fx.As(new(dal.CommentCommand))),
		fx.Annotate(query.NewVideoCommand, fx.As(new(dal.VideoCommand))),
		fx.Annotate(query.NewLikeCommand, fx.As(new(dal.LikeCommand))),
	}
	fx.New(
		sub.Module,
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Provide(
			config.NewConfig,
			zap.NewProduction,
			otel.TracerProvider,
		),
		fx.Provide(dataLayer...),
	).Run()
}
