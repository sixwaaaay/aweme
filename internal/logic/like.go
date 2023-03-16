package logic

import (
	"context"
	"github.com/PlanVX/aweme/internal/dal"
	"github.com/PlanVX/aweme/internal/sub"
	"github.com/PlanVX/aweme/internal/types"
	dapr "github.com/dapr/go-sdk/client"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"time"
)

type (
	// Like is the like logic layer struct
	Like struct {
		likeCommand dal.LikeCommand
		client      dapr.Client
	}
	// LikeParam is the param for NewLike
	LikeParam struct {
		fx.In
		LikeCommand dal.LikeCommand
		Client      dapr.Client
	}
)

// NewLike returns a new Like logic
func NewLike(param LikeParam) *Like {
	return &Like{
		likeCommand: param.LikeCommand,
		client:      param.Client,
	}
}

// Like is the like logic
// handle the like action
func (l *Like) Like(c context.Context, req *types.FavoriteActionReq) (*types.FavoriteActionResp, error) {
	owner, _ := c.Value(ContextKey).(int64) // get the owner from context
	switch req.ActionType {
	case int32(1): // means add like for a video

		like := &dal.Like{
			VideoID: req.VideoID,
			UserID:  owner,
		}

		err := l.likeCommand.Insert(c, like)
		if err != nil {
			return nil, err
		}
	case int32(2): // means remove like for a video

		err := l.likeCommand.Delete(c, req.VideoID, owner)
		if err != nil {
			return nil, err
		}
	default:
		return nil, echo.ErrBadRequest
	}
	return &types.FavoriteActionResp{}, nil
}

// LikeEvent publish the like event to the pub/sub
func (l *Like) LikeEvent(c context.Context, req *types.FavoriteActionReq) (*types.FavoriteActionResp, error) {
	owner, _ := c.Value(ContextKey).(int64) // get the owner from context
	switch req.ActionType {
	case int32(1), int32(2):
		like := &sub.LikeAction{
			VideoID:        req.VideoID,
			UserID:         owner,
			ActionType:     req.ActionType,
			MilliTimestamp: time.Now().UnixMilli(),
		}
		data, err := jsoniter.Marshal(like)
		if err != nil {
			return nil, err
		}
		err = l.client.PublishEvent(c, sub.PubSubName, sub.LikeActionEvent, data)
		if err != nil {
			return nil, err
		}
	default:
		return nil, echo.ErrBadRequest
	}
	return &types.FavoriteActionResp{}, nil
}
