package sub

import "go.uber.org/fx"

var Module = fx.Module("sub",
	fx.Provide(
		fx.Annotate(newCommentEventHandler, fx.As(new(eventHandler)), fx.ResultTags(`group:"evh"`)),
		fx.Annotate(newLikeEventHandler, fx.As(new(eventHandler)), fx.ResultTags(`group:"evh"`)),
		fx.Annotate(newVideoEventHandler, fx.As(new(eventHandler)), fx.ResultTags(`group:"evh"`)),
	),
	fx.Invoke(
		NewDaprService,
	),
)
