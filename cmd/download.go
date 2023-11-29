package cmd

import (
	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"

	"github.com/adamgoose/hlsdl/lib"
	"github.com/adamgoose/hlsdl/lib/jobs"
)

var downloadCmd = &cobra.Command{
	Use:     "download {url} [filename]",
	Aliases: []string{"dl", "d"},
	Args:    cobra.RangeArgs(1, 2),
	RunE: lib.RunE(func(args []string, a *asynq.Client) error {
		if len(args) == 1 {
			args = append(args, "")
		}

		j := jobs.DownloadHLS{
			URL:      args[0],
			FileName: args[1],
		}
		t := j.Task()

		_, err := a.Enqueue(t)
		return err
	}),
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
