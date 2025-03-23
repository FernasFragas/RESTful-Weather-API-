package api

import (
	"context"
	"net/http"
	"weatherservice"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	NumberOfMaxVideos = 3
)

type YoutubeAPI struct {
	client *http.Client

	key string
}

func NewYoutubeAPI(key string) *YoutubeAPI {
	return &YoutubeAPI{
		client: http.DefaultClient,
		key:    key,
	}
}

func (api *YoutubeAPI) FetchReportData(ctx context.Context, query string) (*weatherservice.DataToReport[weatherservice.VideosStream], error) {
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(api.key))
	if err != nil {
		return nil, err
	}

	call := youtubeService.Search.List([]string{"snippet"}).Q(query).MaxResults(NumberOfMaxVideos)

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	var videoStreams weatherservice.VideosStream
	for _, item := range response.Items {
		videoStreams = append(videoStreams, weatherservice.GeneralVideoStreamInfo{
			Title:   item.Snippet.Title,
			VideoID: item.Id.VideoId,
		})
	}

	return &weatherservice.DataToReport[weatherservice.VideosStream]{
		Data: videoStreams,
	}, nil
}

func (api *YoutubeAPI) FetchGeneralInfo(ctx context.Context, _ string) (*weatherservice.DataToReport[weatherservice.VideosStream], error) {
	return nil, nil
}
