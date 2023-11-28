package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "hlsdl",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	viper.SetEnvPrefix("hlsdl")
	viper.AutomaticEnv()

	viper.SetDefault("redis_addr", "localhost:6379")
	viper.SetDefault("redis_db", 0)
}
