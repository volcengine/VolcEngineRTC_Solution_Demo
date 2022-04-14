package task

import (
	"github.com/robfig/cron/v3"
)

var c *cron.Cron

func GetCronTask() *cron.Cron {
	if c == nil {
		c = cron.New()
	}
	return c
}
