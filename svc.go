package app

import (
	"github.com/gopherslab/redbook/internal/cache"
	"github.com/gopherslab/redbook/internal/config"
	emailService "github.com/gopherslab/redbook/internal/notifications/email"
	emailModel "github.com/gopherslab/redbook/internal/notifications/email/model"
	authService "github.com/gopherslab/redbook/internal/service/auth"
	authModel "github.com/gopherslab/redbook/internal/service/auth/model"
	authRepo "github.com/gopherslab/redbook/internal/service/auth/repo"
	beachService "github.com/gopherslab/redbook/internal/service/beach"
	beachModel "github.com/gopherslab/redbook/internal/service/beach/model"
	beachRepo "github.com/gopherslab/redbook/internal/service/beach/repo"
	clientService "github.com/gopherslab/redbook/internal/service/client"
	clientModel "github.com/gopherslab/redbook/internal/service/client/model"
	clientRepo "github.com/gopherslab/redbook/internal/service/client/repo"
	employeeService "github.com/gopherslab/redbook/internal/service/employee"
	employeeModel "github.com/gopherslab/redbook/internal/service/employee/model"
	employeeRepo "github.com/gopherslab/redbook/internal/service/employee/repo"
	metadataService "github.com/gopherslab/redbook/internal/service/metadata"
	metadataModel "github.com/gopherslab/redbook/internal/service/metadata/model"
	metadataRepo "github.com/gopherslab/redbook/internal/service/metadata/repo"
	needService "github.com/gopherslab/redbook/internal/service/need"
	needModel "github.com/gopherslab/redbook/internal/service/need/model"
	needRepo "github.com/gopherslab/redbook/internal/service/need/repo"
	projectService "github.com/gopherslab/redbook/internal/service/project"
	projectModel "github.com/gopherslab/redbook/internal/service/project/model"
	projectRepo "github.com/gopherslab/redbook/internal/service/project/repo"
	reportService "github.com/gopherslab/redbook/internal/service/reports"
	taggingService "github.com/gopherslab/redbook/internal/service/tagging"
	taggingModel "github.com/gopherslab/redbook/internal/service/tagging/model"
	taggingRepo "github.com/gopherslab/redbook/internal/service/tagging/repo"

	reportModel "github.com/gopherslab/redbook/internal/service/reports/model"
	reportRepo "github.com/gopherslab/redbook/internal/service/reports/repo"

	"gorm.io/gorm"
)

type services struct {
	authSvc     authModel.Service
	clientSvc   clientModel.Service
	projectSvc  projectModel.Service
	employeeSvc employeeModel.Service
	beachSvc    beachModel.Service
	needSvc     needModel.Service
	taggingSvc  taggingModel.Service
	metadataSvc metadataModel.Service
	mailSvc     emailModel.Service
	reportSvc   reportModel.Service
}

type repos struct {
	authRepo     authModel.Repository
	clientRepo   clientModel.Repository
	projectRepo  projectModel.Repository
	employeeRepo employeeModel.Repository
	beachRepo    beachModel.Repository
	needRepo     needModel.Repository
	taggingRepo  taggingModel.Repository
	metadataRepo metadataModel.Repository
	reportRepo   reportRepo.Repository
}

func buildServices(
	cfg *config.Config,
	db *gorm.DB,
	cache cache.Cache,
) *services {
	svc := &services{}
	repo := &repos{}
	repo.buildRepos(db)
	svc.buildAuthService(repo, &cfg.Auth, cache)
	svc.buildBeachService(repo)
	svc.buildClientService(repo)
	svc.buildEmployeeService(repo, cfg)
	svc.buildneedService(repo)
	svc.buildProjectService(repo)
	svc.buildTaggingService(repo)
	svc.buildMetadataService(repo)
	svc.buildEmailService(cfg)
	svc.buildReportService(repo)

	return svc
}

func (r *repos) buildRepos(db *gorm.DB) {
	r.authRepo = authRepo.NewRepository(db)
	r.clientRepo = clientRepo.NewRepository(db)
	r.projectRepo = projectRepo.NewRepository(db)
	r.employeeRepo = employeeRepo.NewRepository(db)
	r.beachRepo = beachRepo.NewRepository(db)
	r.needRepo = needRepo.NewRepository(db)
	r.taggingRepo = taggingRepo.NewRepository(db)
	r.metadataRepo = metadataRepo.NewRepository(db)
	r.reportRepo = *reportRepo.NewRepository(db)
}

func (s *services) buildAuthService(repo *repos, cfg *authModel.Config, cache cache.Cache) {
	s.authSvc = authService.NewService(cfg, repo.authRepo, cache)
}

func (s *services) buildClientService(repo *repos) {
	s.clientSvc = clientService.NewService(repo.clientRepo, repo.projectRepo)
}

func (s *services) buildProjectService(repo *repos) {
	s.projectSvc = projectService.NewService(repo.projectRepo, repo.clientRepo, repo.taggingRepo, repo.needRepo)
}

func (s *services) buildEmployeeService(repo *repos, cfg *config.Config) {
	s.employeeSvc = employeeService.NewService(repo.employeeRepo, repo.taggingRepo, repo.beachRepo, repo.metadataRepo, cfg)
}

func (s *services) buildBeachService(repo *repos) {
	s.beachSvc = beachService.NewService(repo.beachRepo, repo.taggingRepo, repo.employeeRepo, repo.needRepo)
}

func (s *services) buildneedService(repo *repos) {
	s.needSvc = needService.NewService(repo.needRepo, repo.projectRepo, repo.employeeRepo)
}

func (s *services) buildTaggingService(repo *repos) {
	s.taggingSvc = taggingService.NewService(repo.taggingRepo, repo.projectRepo, repo.employeeRepo)
}

func (s *services) buildMetadataService(repo *repos) {
	s.metadataSvc = metadataService.NewService(repo.metadataRepo)
}

func (s *services) buildEmailService(cfg *config.Config) {
	s.mailSvc = emailService.NewService(&cfg.Email)
}

func (s *services) buildReportService(repo *repos) {
	s.reportSvc = reportService.NewService(&repo.reportRepo)
}
