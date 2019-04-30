package managers

import (
	"frank/src/go/config"
	"frank/src/go/helpers/log"
	"frank/src/go/models"
	"gopkg.in/robfig/cron.v2"
)

const DEFAULT_DDNS_CRON = "0/2 * * * *"

var c = cron.New()

func AddCron(cron models.Cron) {
	action, err := config.GetAction(cron.Action)
	if err != nil {
		log.Log.Errorf("Cannot find Action %s", cron.Action)
		return
	}

	c.AddFunc(cron.CronExpression, func() {
		log.Log.Infof("Cron Started '%s'", cron.Description)
		ActivePlugins.ExecAction(action, cron.Extra)
	})
}

func LoadCrons() {
	c.Start()
	c.Stop()
	c = cron.New()

	crons := config.ParsedConfig.Crons

	for _, c := range crons {
		AddCron(c)
	}

	if config.ParsedConfig.Ddns.Hostname != "" {
		cronExpression := DEFAULT_DDNS_CRON
		if config.ParsedConfig.Ddns.CronExpression != "" {
			cronExpression = config.ParsedConfig.Ddns.CronExpression
		}

		c.AddFunc(cronExpression, func() {
			log.Log.Infof("Cron Started 'DDNS'")
			go DdnsManager.SetIp()
		})

	}

	c.Start()
}
