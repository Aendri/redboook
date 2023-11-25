package app

import (
	"context"
	"time"

	"github.com/gopherslab/redbook/internal/config"
	jobModel "github.com/gopherslab/redbook/internal/jobs"
	beachJob "github.com/gopherslab/redbook/internal/jobs/beach"
	"github.com/gopherslab/redbook/internal/jobs/deallocation"
	needsJob "github.com/gopherslab/redbook/internal/jobs/needs"
	zohoJob "github.com/gopherslab/redbook/internal/jobs/zoho"

	"github.com/gopherslab/redbook/pkg/log"

	"github.com/robfig/cron/v3"
)

type jobs []jobModel.Job

func buildJobs(
	_ log.Logger,
	cfg *config.Config,
	services *services,
) jobs {
	jbs := make(jobs, 0)
	zjob := zohoJob.NewJob(services.employeeSvc, &cfg.Jobs.Zoho)
	needsJob := needsJob.NewJob(services.mailSvc, services.beachSvc, services.needSvc, &cfg.Jobs.NeedNBeach)
	beachJob := beachJob.NewJob(services.employeeSvc, services.taggingSvc, services.projectSvc, &cfg.Jobs.Beach)
	deallocationJob := deallocation.NewJob(services.mailSvc, services.taggingSvc, &cfg.Jobs.DeallocationList)
	jbs = append(jbs, zjob, needsJob, beachJob, deallocationJob)
	return jbs
}

func initCron(
	ctx context.Context,
	log log.Logger,
	jbs jobs,
) *cron.Cron {
	ind, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Error("failed to load location", err)
	}
	c := cron.New(cron.WithLocation(ind))
	for _, job := range jbs {
		_, err = c.AddFunc(job.GetCron(), job.GetJob(ctx, log))
		if err != nil {
			log.Error("adding job", job.GetName(), err)
		}
	}
	return c
}

func (a *Application) startJobs() {
	a.cron.Start()
}

func (a *Application) stopJobs() {
	a.cron.Stop()
}
