package managers

import (
	"frank/src/go/config"
	"frank/src/go/helpers/log"
	"frank/src/go/models"

	"github.com/jasonlvhit/gocron"
)

func actionCronJob(action models.Action, extra map[string]string) {
	log.Log.Warning("first Job")
	ActivePlugins.ExecAction(action, extra)
}

func AddCron(cron models.Cron) {
	action, err := config.GetAction(cron.Action)
	if err != nil {
		log.Log.Errorf("Cannot find Action %s", cron.Action)
		return
	}

	gc := gocron.Every(uint64(cron.Every))
	if cron.TimeType == "seconds" {
		gc = gc.Seconds()
	}

	gc.Do(actionCronJob, action, cron.Extra)
}

func LoadCrons() {
	gocron.Clear()
	crons := config.ParsedConfig.Crons

	for _, c := range crons {
		AddCron(c)
	}

	gocron.Start()
}
