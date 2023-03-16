package logic

import (
	"context"
	"github.com/PlanVX/aweme/internal/dal"
	"github.com/PlanVX/aweme/internal/sub"
	"github.com/PlanVX/aweme/internal/types"
	dapr "github.com/dapr/go-sdk/client"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/fx"
	"strconv"
	"time"
)

type (
	// Upload is the logic for upload video
	Upload struct {
		videoCommand dal.VideoCommand
		client       dapr.Client
		uploader     Uploader
	}
	// UploadParam is the param for NewUpload
	UploadParam struct {
		fx.In
		VideoCommand dal.VideoCommand
		Client       dapr.Client
		Helper       Uploader
	}
)

// NewUpload returns a new Upload logic
func NewUpload(param UploadParam) *Upload {
	return &Upload{
		videoCommand: param.VideoCommand,
		client:       param.Client,
		uploader:     param.Helper,
	}
}

// UploadVideo publishes a video
func (u *Upload) UploadVideo(c context.Context, req *types.UploadReq) (*types.UploadResp, error) {
	owner, _ := c.Value(ContextKey).(int64) // get the owner from context

	key := strconv.FormatInt(time.Now().UnixNano(), 10) + req.FileName // generate a unique key for the video

	upload, err := u.uploader.Upload(UploadInput{Key: key, Value: req.Data}) // upload the video
	if err != nil {
		return nil, err
	}

	// insert the video into database
	err = u.videoCommand.Insert(c, &dal.Video{VideoURL: upload, UserID: owner, Title: req.Title})
	if err != nil {
		return nil, err
	}

	return &types.UploadResp{}, nil
}

// UploadVideoEvent is the event handler for video upload
func (u *Upload) UploadVideoEvent(c context.Context, req *types.UploadReq) (*types.UploadResp, error) {
	owner, _ := c.Value(ContextKey).(int64) // get the owner from context

	key := strconv.FormatInt(time.Now().UnixNano(), 10) + req.FileName // generate a unique key for the video

	upload, err := u.uploader.Upload(UploadInput{Key: key, Value: req.Data}) // upload the video
	if err != nil {
		return nil, err
	}

	// insert the video into database
	s := &sub.Video{
		VideoURL:       upload,
		UserID:         owner,
		Title:          req.Title,
		MilliTimestamp: time.Now().UnixMilli(),
	}
	data, err := jsoniter.Marshal(s)
	if err != nil {
		return nil, err
	}
	err = u.client.PublishEvent(c, sub.PubSubName, sub.VideoEvent, data)
	if err != nil {
		return nil, err
	}

	return &types.UploadResp{}, nil
}
