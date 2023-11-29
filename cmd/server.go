package cmd

import (
	"os"

	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/adamgoose/hlsdl/lib"
	"github.com/adamgoose/hlsdl/lib/jobs"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"s"},
	RunE: lib.RunE(func(rco asynq.RedisClientOpt) error {
		if err := os.MkdirAll(viper.GetString("out"), 0755); err != nil {
			return err
		}

		srv := asynq.NewServer(rco,
			asynq.Config{
				Concurrency: viper.GetInt("concurrency"),
				Queues: map[string]int{
					"critical": 6,
					"default":  3,
					"low":      1,
				},
			},
		)
		mux := asynq.NewServeMux()

		mux.Handle(jobs.DownloadHLSJob, &jobs.DownloadHLSHandler{})

		return srv.Run(mux)
	}),
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntP("concurrency", "c", 5, "Number of concurrent workers")
	viper.BindPFlag("concurrency", serverCmd.Flags().Lookup("concurrency"))

	serverCmd.Flags().StringP("out", "o", "output", "Output directory")
	viper.BindPFlag("out", serverCmd.Flags().Lookup("out"))
}
