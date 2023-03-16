package sub

import (
	"context"
	"github.com/PlanVX/aweme/internal/dal"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/go-sdk/service/grpc"
	"github.com/go-sql-driver/mysql"
	"github.com/json-iterator/go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"time"
)

type (
	CommentAction struct {
		ID             int64  `json:"id"`
		VideoID        int64  `json:"video_id"`
		Comment        string `json:"comment"`
		UserID         int64  `json:"user_id"`
		ActionType     int8   `json:"action_type"`
		MilliTimestamp int64  `json:"timestamp"` // milliseconds timestamp
	}

	LikeAction struct {
		VideoID        int64 `json:"video_id"`
		UserID         int64 `json:"user_id"`
		ActionType     int32 `json:"action_type"`
		MilliTimestamp int64 `json:"timestamp"` // milliseconds timestamp
	}

	Video struct {
		VideoURL       string `json:"video_url"`
		UserID         int64  `json:"user_id"`
		Title          string `json:"title"`
		MilliTimestamp int64  `json:"timestamp"` // milliseconds timestamp
	}
)

const (
	LikeActionEvent    = "like"
	CommentActionEvent = "comment"
	VideoEvent         = "video"
	PubSubName         = "redis-pubsub"
)

type eventHandler interface {
	subscription() *common.Subscription
	Handle(ctx context.Context, e *common.TopicEvent) (retry bool, err error)
}

type Option struct {
	fx.In
	Handlers []eventHandler `group:"evh"`
	Lf       fx.Lifecycle
}

func NewDaprService(option Option) (common.Service, error) {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "6010"
	}

	s, err := grpc.NewService(":" + appPort)
	if err != nil {
		return nil, err
	}

	for _, handler := range option.Handlers {
		err := s.AddTopicEventHandler(handler.subscription(), handler.Handle)
		if err != nil {
			return nil, err
		}
	}
	err = s.Start()
	if err != nil {
		return nil, err
	}
	option.Lf.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				return s.GracefulStop()
			},
		},
	)
	return s, nil
}

type commentEventHandler struct {
	comment dal.CommentCommand
	logger  *zap.Logger
}

func newCommentEventHandler(comment dal.CommentCommand, logger *zap.Logger) *commentEventHandler {
	return &commentEventHandler{
		comment: comment,
		logger:  logger,
	}
}

func (c *commentEventHandler) Handle(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	v := new(CommentAction)
	err = jsoniter.Unmarshal(e.RawData, v)
	if err != nil {
		return true, err
	}
	var comment dal.Comment
	comment.ID = v.ID
	comment.VideoID = v.VideoID
	comment.UserID = v.UserID
	comment.Content = v.Comment
	if v.ActionType == 1 {
		err = c.comment.Insert(ctx, &comment)
	} else {
		err = c.comment.Delete(ctx, comment.ID, comment.UserID, comment.VideoID)
	}
	if err != nil {

	}
	return false, nil
}

func (c *commentEventHandler) subscription() *common.Subscription {
	return &common.Subscription{
		PubsubName: "redis-pubsub",
		Topic:      "comment",
		Route:      "/comment",
	}
}

type likeEventHandler struct {
	like   dal.LikeCommand
	logger *zap.Logger
}

func newLikeEventHandler(like dal.LikeCommand, logger *zap.Logger) *likeEventHandler {
	return &likeEventHandler{
		like:   like,
		logger: logger,
	}
}

func (l *likeEventHandler) Handle(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	v := new(LikeAction)
	err = jsoniter.Unmarshal(e.RawData, v)
	if err != nil {
		l.logger.Error("cannot unmarshal data", zap.String("event_id", e.ID))
		return false, err
	}
	var like dal.Like
	like.VideoID = v.VideoID
	like.UserID = v.UserID
	if v.ActionType == 1 {
		err = l.like.Insert(ctx, &like)
	} else {
		err = l.like.Delete(ctx, like.UserID, like.VideoID)
	}
	if err != nil {
		retry = handlerMySQLError(err)
		l.logger.Error("event handle error", zap.String("event_id", e.ID), zap.Error(err))
		return retry, err
	}

	l.logger.Info("event handled", zap.String("event_id", e.ID), zap.String("event_topic", e.Topic))
	return false, nil
}

func (l *likeEventHandler) subscription() *common.Subscription {
	return &common.Subscription{
		PubsubName: "redis-pubsub",
		Topic:      "like",
		Route:      "/like",
	}
}

type videoEventHandler struct {
	video  dal.VideoCommand
	logger *zap.Logger
}

func newVideoEventHandler(video dal.VideoCommand, logger *zap.Logger) *videoEventHandler {
	return &videoEventHandler{
		video:  video,
		logger: logger,
	}
}

func (v *videoEventHandler) Handle(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	value := new(Video)
	err = jsoniter.Unmarshal(e.RawData, value)
	if err != nil {
		v.logger.Error("cannot unmarshal data", zap.String("event_id", e.ID))
		return false, err
	}
	var video dal.Video
	video.UserID = value.UserID
	video.VideoURL = value.VideoURL
	video.Title = value.Title
	video.CreatedAt = time.Unix(0, value.MilliTimestamp*int64(time.Millisecond))
	err = v.video.Insert(ctx, &video)
	if err != nil {
		retry = handlerMySQLError(err)
		v.logger.Error("cannot insert video", zap.Error(err), zap.String("event_id", e.ID))
		return retry, err
	}
	v.logger.Info("event handled", zap.String("event_id", e.ID), zap.String("event_topic", e.Topic))
	return false, nil
}

func (v *videoEventHandler) subscription() *common.Subscription {
	return &common.Subscription{
		PubsubName: "redis-pubsub",
		Topic:      "video",
		Route:      "/video",
	}
}

func handlerMySQLError(err error) bool {
	if errMySQL, ok := err.(*mysql.MySQLError); ok {
		switch errMySQL.Number {
		case 1062:
			return false
		}
	}
	return true
}
