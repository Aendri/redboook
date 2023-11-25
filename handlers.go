package app

import (
	"github.com/gopherslab/redbook/api"
	"github.com/gopherslab/redbook/api/v1/auth"
	"github.com/gopherslab/redbook/api/v1/beach"
	"github.com/gopherslab/redbook/api/v1/client"
	"github.com/gopherslab/redbook/api/v1/employee"
	"github.com/gopherslab/redbook/api/v1/healthcheck"
	"github.com/gopherslab/redbook/api/v1/metadata"
	"github.com/gopherslab/redbook/api/v1/need"
	"github.com/gopherslab/redbook/api/v1/project"
	"github.com/gopherslab/redbook/api/v1/reports"
	"github.com/gopherslab/redbook/api/v1/tagging"

	authService "github.com/gopherslab/redbook/internal/service/auth"
)

func (a *Application) SetupHandlers() {
	authMiddleware := authService.NewAuthMiddleWare(a.cache, a.log, a.services.authSvc, a.cfg.Enabled.Auth)
	validator := api.NewValidations(a.services.metadataSvc)
	beach.RegisterHandlers(
		a.router,
		a.services.beachSvc,
		a.log,
		authMiddleware.Auth,
	)
	employee.RegisterHandlers(
		a.router,
		a.services.employeeSvc,
		a.log,
		authMiddleware.Auth,
		validator,
	)
	need.RegisterHandlers(
		a.router,
		a.services.needSvc,
		a.services.taggingSvc,
		a.log,
		authMiddleware.Auth,
		validator,
	)
	project.RegisterHandlers(
		a.router,
		a.services.projectSvc,
		a.log,
		authMiddleware.Auth,
		validator,
	)

	tagging.RegisterHandlers(
		a.router,
		a.services.taggingSvc,
		a.log,
		authMiddleware.Auth,
		validator,
	)
	client.RegisterHandlers(
		a.router,
		a.services.clientSvc,
		a.services.projectSvc,
		a.log,
		authMiddleware.Auth,
		validator,
	)

	auth.RegisterHandlers(
		a.router,
		a.services.authSvc,
		a.log,
		authMiddleware.Auth,
	)

	metadata.RegisterHandlers(
		a.router,
		a.services.metadataSvc,
		a.log,
		authMiddleware.Auth,
		validator,
	)
	reports.RegisterHandlers(a.router,
		a.services.reportSvc,
		a.log,
		authMiddleware.Auth,
		validator,
	)

	healthcheck.RegisterHandlers(a.router)
}
