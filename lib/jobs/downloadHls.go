package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const DownloadHLSJob = "download-hls"

type DownloadHLS struct {
	URL      string
	FileName string
}

func (p *DownloadHLS) Task() *asynq.Task {
	j, _ := json.Marshal(p)
	return asynq.NewTask(DownloadHLSJob, j)
}

func (p *DownloadHLS) Unmarshal(t *asynq.Task) error {
	return json.Unmarshal(t.Payload(), p)
}

type DownloadHLSHandler struct {
	//
}

func (h *DownloadHLSHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p DownloadHLS
	if err := p.Unmarshal(t); err != nil {
		return err
	}

	if p.FileName == "" {
		p.FileName = time.Now().Format("2006-01-02_15-04-05")
	}

	out := fmt.Sprintf("output/%s.mp4", p.FileName)
	err := ffmpeg.Input(p.URL, nil).Output(out, ffmpeg.KwArgs{"c:v": "copy", "c:a": "copy"}).Run()
	if err != nil {
		log.Println(err)
	}

	return err
}
