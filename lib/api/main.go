package api

import (
	"net/http"
	"sync"

	"github.com/adamgoose/hlsdl/lib/jobs"
	"github.com/defval/di"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type API struct {
	di.Inject
	Asynq *asynq.Client
}

func (a *API) Run(wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "HLSDL")
	})

	e.GET("/dl", func(c echo.Context) error {
		url := c.QueryParam("url")
		filename := c.QueryParam("filename")

		if url == "" {
			return c.String(http.StatusBadRequest, "URL is required")
		}

		j := jobs.DownloadHLS{
			URL:      url,
			FileName: filename,
		}
		t := j.Task()

		_, err := a.Asynq.Enqueue(t)
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, "Queued")
	})

	errChan <- e.Start(viper.GetString("listen"))
}
