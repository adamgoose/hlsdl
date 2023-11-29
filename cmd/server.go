package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/adamgoose/hlsdl/lib"
	"github.com/adamgoose/hlsdl/lib/api"
	"github.com/adamgoose/hlsdl/lib/jobs"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"s"},
	RunE: lib.RunE(func(js *jobs.Server, as *api.API) error {
		if err := os.MkdirAll(viper.GetString("out"), 0755); err != nil {
			return err
		}

		var wg sync.WaitGroup
		jobErrChan := make(chan error, 1)
		apiErrChan := make(chan error, 1)

		wg.Add(1)
		go js.Run(&wg, jobErrChan)

		wg.Add(1)
		go as.Run(&wg, apiErrChan)

		select {
		case err := <-jobErrChan:
			fmt.Println("Asynq server error:", err)
		case err := <-apiErrChan:
			fmt.Println("API server error:", err)
		}

		return nil
	}),
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("listen", "l", ":8881", "Listen address for HTTP")
	viper.BindPFlag("listen", serverCmd.Flags().Lookup("listen"))

	serverCmd.Flags().IntP("concurrency", "c", 5, "Number of concurrent workers")
	viper.BindPFlag("concurrency", serverCmd.Flags().Lookup("concurrency"))

	serverCmd.Flags().StringP("out", "o", "output", "Output directory")
	viper.BindPFlag("out", serverCmd.Flags().Lookup("out"))
}
