package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	"github.com/new-york-services/api_gateway/api"
	"github.com/new-york-services/api_gateway/config"
	"github.com/new-york-services/api_gateway/pkg/logger"
	"github.com/new-york-services/api_gateway/services"
	"github.com/new-york-services/api_gateway/storage/redis"

	r "github.com/gomodule/redigo/redis"

	gormadapter "github.com/casbin/gorm-adapter/v2"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.String("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	psqlString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostServiceHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase)

	a, err := gormadapter.NewAdapter("postgres", psqlString, true)
	if err != nil {
		log.Error("new adapter error: ", logger.Error(err))
		return
	}

	casbinEnForcer, err := casbin.NewEnforcer(cfg.CasbinConfigPath, a)
	if err != nil {
		log.Error("new enforcer error", logger.Error(err))
		return
	}

	err = casbinEnForcer.LoadPolicy()
	if err != nil {
		log.Error("casbin load policy error", logger.Error(err))
	}

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error: ", logger.Error(err))
	}

	pool := r.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (r.Conn, error) {
			c, err := r.Dial("tcp", fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	server := api.New(api.Option{
		Conf:            cfg,
		ServiceManager:  serviceManager,
		Logger:          log,
		InMemoryStorage: redis.NewRedisRepo(&pool),
		CasbinEnforcer:  casbinEnForcer,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run HTTP server: ", logger.Error(err))
		panic(err)
	}

	casbinEnForcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddDomainMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnForcer.GetRoleManager().(*defaultrolemanager.RoleManager).AddDomainMatchingFunc("keyMatch3", util.KeyMatch3)
}
