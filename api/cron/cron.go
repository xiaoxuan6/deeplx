package cron

import (
	"github.com/robfig/cron/v3"
	"github.com/xiaoxuan6/deeplx"
)

func Start() {
	c := cron.New(cron.WithSeconds())
	_, _ = c.AddFunc("0 0 23 * * *", func() {
		deeplx.CheckUrlAndReloadBlack()
		deeplx.LoadBlack(true)
	})

	c.Start()
}
