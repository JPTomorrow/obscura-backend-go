package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"github.com/JPTomorrow/obscura/db"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YTService struct {
	*youtube.Service
	videoPool *YTVideoPool
}

type YTVideoPool struct {
	videos []*db.YoutubeVideo
	mu     sync.Mutex
}

func NewYoutubeService() *YTService {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	return &YTService{
		service,
		&YTVideoPool{},
	}
}

// pull
func (yt *YTService) PullNewVideos(limit int64) error {
	// now := time.Now()
	// pbTime, err := time.Parse(time.RFC3339, now.String())
	// if err != nil {
	// 	return err
	// }

	s, err := yt.Search.List([]string{"id,snippet"}).VideoDuration("medium").VideoDimension("2d").VideoEmbeddable("true").RelevanceLanguage("en").Type("video").MaxResults(limit).Q("horror|ghost|apparition|hunter|spectre|demon|haunt|spooky|camera|creepy|scary").SafeSearch("none").Do()
	if err != nil {
		return err
	}

	yt.videoPool.mu.Lock()
	for _, item := range s.Items {
		if item.Snippet.Description == "" {
			continue
		}

		newVid := &db.YoutubeVideo{
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			VideoTag:    item.Id.VideoId,
		}
		yt.videoPool.videos = append(yt.videoPool.videos, newVid)
	}
	yt.videoPool.mu.Unlock()

	// print video pool
	// for _, vid := range videoPool.videos {
	// 	fmt.Println(vid)
	// }

	return nil
}

func (yt *YTService) PitchRandomVideo() (*db.YoutubeVideo, error) {
	if len(yt.videoPool.videos) == 0 {
		return nil, fmt.Errorf("no videos in pool")
	}
	rng := rand.Intn(len(yt.videoPool.videos)-0) + 0
	return yt.videoPool.videos[rng], nil
}
