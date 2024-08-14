package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/adamgoose/hlsdl/lib/jobs"
	"github.com/defval/di"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type API struct {
	di.Inject
	Asynq    *asynq.Client
	Asynqmon *asynqmon.HTTPHandler
}

func (a *API) Run(wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "HLSDL")
	})

	e.Any("/asynq*", echo.WrapHandler(a.Asynqmon))

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

		_, err := a.Asynq.Enqueue(t, asynq.Retention(24*time.Hour))
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, "Queued")
	})

	errChan <- e.Start(viper.GetString("listen"))
}
