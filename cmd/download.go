package cmd

import (
	"github.com/adamgoose/hlsdl/lib"
	"github.com/adamgoose/hlsdl/lib/jobs"
	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:     "download",
	Aliases: []string{"dl", "d"},
	RunE: lib.RunE(func(args []string, a *asynq.Client) error {
		j := jobs.DownloadHLS{
			URL: args[0],
		}
		t := j.Task()

		_, err := a.Enqueue(t)
		return err
	}),
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
