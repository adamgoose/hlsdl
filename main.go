package main

import (
	"github.com/adamgoose/hlsdl/cmd"
	"github.com/adamgoose/hlsdl/lib"
	"github.com/defval/di"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

func main() {
	if err := lib.Apply(
		di.Provide(func() asynq.RedisClientOpt {
			return asynq.RedisClientOpt{
				Addr:     viper.GetString("redis_addr"),
				Password: viper.GetString("redis_pass"),
				DB:       viper.GetInt("redis_db"),
			}
		}),
		di.Provide(func(r asynq.RedisClientOpt) *asynq.Client {
			return asynq.NewClient(r)
		}),
	); err != nil {
		panic(err)
	}

	cmd.Execute()
}
