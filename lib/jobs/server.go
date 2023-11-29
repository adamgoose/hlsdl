package jobs

import (
	"sync"

	"github.com/defval/di"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

type Server struct {
	di.Inject
	RedisClientOpt asynq.RedisClientOpt
}

func (j *Server) Run(wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	srv := asynq.NewServer(j.RedisClientOpt, asynq.Config{
		Concurrency: viper.GetInt("concurrency"),
	})
	mux := asynq.NewServeMux()

	mux.Handle(DownloadHLSJob, &DownloadHLSHandler{})

	errChan <- srv.Run(mux)
}
