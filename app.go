package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"

	mdrpc "github.com/gopherslab/redbook/api/rpc/metadata"
	rpc "github.com/gopherslab/redbook/api/rpc/redbook"
	"github.com/gopherslab/redbook/cmd/redbook/migrations"
	"github.com/gopherslab/redbook/internal/cache"
	"github.com/gopherslab/redbook/internal/cache/redis"
	"github.com/gopherslab/redbook/internal/config"

	"github.com/gopherslab/redbook/internal/pb/redbook"
	metadataModel "github.com/gopherslab/redbook/internal/service/metadata/model"
	"github.com/gopherslab/redbook/pkg/db"
	"github.com/gopherslab/redbook/pkg/log"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

const timeout = 5 * time.Second

type Application struct {
	cache      cache.Cache
	db         *gorm.DB
	log        log.Logger
	cfg        *config.Config
	router     *mux.Router
	httpServer *http.Server
	grpcServer *grpc.Server
	services   *services
	engines    engines
	jobs       jobs
	cron       *cron.Cron
}

func (a *Application) Init(ctx context.Context, configFile string, migrationPath string, seedDataPath string) {
	log := log.New().With(ctx)
	a.log = log

	config, err := config.Load(log, configFile)
	if err != nil {
		log.Fatalf("failed to read config: %s ", err)
		return
	}
	a.cfg = config

	zerolog.SetGlobalLevel(zerolog.Level(config.Logging.Level))

	db, err := db.NewDB(config.DBConfig, log)
	if err != nil {
		log.Fatalf("error connecting db: %s ", err)
		return
	}
	migrations.Migrate(log, db, config.DBConfig, migrationPath)
	a.db = db

	cache, err := redis.NewPoolCache(config.Redis, log)
	if err != nil {
		log.Fatalf("failed to create redis client. continuing to initialize app: %w", err)
	}
	a.cache = cache

	router := mux.NewRouter()
	a.router = router

	services := buildServices(config, db, cache)
	a.services = services
	a.engines = buildEngines(a.log, a.cfg, a.services)
	a.jobs = buildJobs(a.log, a.cfg, a.services)
	a.cron = initCron(ctx, a.log, a.jobs)
	if a.cfg.Enabled.Refresh {
		err = a.services.employeeSvc.Refresh(ctx, a.log)
		if err != nil {
			log.Errorf("failed to refresh employees. continuing to initialize app: %v", err)
		}
	}

	migrations.SeedData(log, db, config.DBConfig, seedDataPath)

	_, err = a.services.metadataSvc.Get(ctx, a.log, metadataModel.Filter{})
	if err != nil {
		log.Errorf("failed to cache metadata. continuing to initialize app: %v", err)
	}
	a.SetupHandlers()
}

func (a *Application) Start(ctx context.Context) {
	a.router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
		AllowedHeaders:   []string{"accept", "Authorization", "content-type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler)

	// start engines
	if engStartErr := a.engines.Start(ctx, a.log); engStartErr != nil {
		a.log.Error("error starting engines")
	}

	a.httpServer = &http.Server{
		Addr:              ":" + fmt.Sprintf("%v", a.cfg.Server.Port),
		Handler:           a.router,
		ReadHeaderTimeout: timeout,
	}
	go func() {
		defer a.log.Infof("server stopped listening")
		if a.cfg.Enabled.HTTP {
			if err := a.httpServer.ListenAndServe(); err != nil {
				a.log.Errorf("failed to listen and serve: %v ", err)
				return
			}
			return
		}
		if err := a.httpServer.ListenAndServeTLS("/etc/letsencrypt/live/rbtest-api.gopherslab.com/fullchain.pem",
			"/etc/letsencrypt/live/rbtest-api.gopherslab.com/privkey.pem"); err != nil {
			a.log.Fatalf("failed to listen and serve: %s ", err)
			return
		}
	}()
	a.log.Infof("http server started on %d ...", a.cfg.Server.Port)

	go func() {
		// grpc server
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.Grpc.Host, a.cfg.Grpc.Port))
		if err != nil {
			a.log.Fatalf("failed to listen: %v", err)
		}

		a.grpcServer = grpc.NewServer()

		redbook.RegisterRedbookServiceServer(
			a.grpcServer,
			rpc.NewServer(
				a.services.employeeSvc,
				a.services.projectSvc,
				a.services.clientSvc,
				a.services.taggingSvc,
				a.log),
		)
		redbook.RegisterMetaDataServiceServer(
			a.grpcServer,
			mdrpc.NewServer(
				a.services.metadataSvc,
				a.log),
		)
		if err := a.grpcServer.Serve(lis); err != nil {
			a.log.Fatalf("failed to serve: %v", err)
		}
	}()

	a.log.Infof("grpc server started on %d ...", a.cfg.Grpc.Port)
	a.startJobs()
}

func (a *Application) Stop(ctx context.Context) {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		a.log.Error(err)
	}
	a.grpcServer.GracefulStop()
	err = a.engines.Stop(ctx, a.log)
	if err != nil {
		a.log.Error(err)
	}
	a.stopJobs()
	a.log.Info("shutting down....")
}
