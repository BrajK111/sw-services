package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/egovernments/sw-services-go/internal/config"
	"github.com/egovernments/sw-services-go/internal/repository/postgres"
	"github.com/egovernments/sw-services-go/internal/service"
	"github.com/egovernments/sw-services-go/internal/transport/http/handler"
	"github.com/egovernments/sw-services-go/internal/transport/http/middleware"
	"github.com/egovernments/sw-services-go/internal/util"
	"github.com/egovernments/sw-services-go/internal/validator"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/application.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgres.Connect(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	repo := postgres.NewConnectionRepository(db)

	propertyClient := util.NewPropertyClient(cfg.Peers.PropertyHost, cfg.Peers.PropertySearchPath)
	mdmsClient := util.NewMdmsClient(cfg.Peers.MdmsHost, cfg.Peers.MdmsSearchPath)
	idgenClient := util.NewIdGenClient(cfg.Peers.IdgenHost, cfg.Peers.IdgenGeneratePath)
	userClient := util.NewUserClient(cfg.Peers.UserHost, cfg.Peers.UserSearchPath, cfg.Peers.UserCreateNoValidate)

	propertyValidator := validator.NewPropertyValidator(propertyClient)
	mdmsValidator := validator.NewMDMSValidator(mdmsClient)

	svc := service.NewConnectionService(repo, cfg.Pagination, propertyValidator, mdmsValidator, idgenClient, userClient)
	connectionHandler := handler.NewConnectionHandler(svc)

	router := gin.Default()
	group := router.Group(cfg.Server.ContextPath)
	{
		// Public — no auth required.
		group.GET("/actuator/health", handler.Health)

		// CITIZEN, SW_CEMP, SUPERUSER: apply for a new connection.
		group.POST("/swc/_create",
			middleware.RequireRole(middleware.RoleCitizen, middleware.RoleSWCEMP, middleware.RoleSuperUser),
			connectionHandler.Create)

		// Employee roles only: advance the workflow (field inspection → approval → active).
		group.POST("/swc/_update",
			middleware.RequireRole(middleware.RoleSWCEMP, middleware.RoleSWFieldInspector, middleware.RoleSWApprover, middleware.RoleSuperUser, middleware.RoleEmployee),
			connectionHandler.Update)

		// All authenticated roles may search.
		group.POST("/swc/_search",
			middleware.RequireRole(middleware.RoleCitizen, middleware.RoleEmployee, middleware.RoleSWCEMP, middleware.RoleSWFieldInspector, middleware.RoleSWApprover, middleware.RoleSuperUser),
			connectionHandler.Search)

		// Internal use only: SYSTEM, SUPERUSER (called by sw-calculator).
		group.POST("/swc/_plainsearch",
			middleware.RequireRole(middleware.RoleSystem, middleware.RoleSuperUser, middleware.RoleSWApprover),
			connectionHandler.PlainSearch)

		group.POST("/swc/_encryptOldData",
			middleware.RequireRole(middleware.RoleSuperUser),
			connectionHandler.EncryptOldData)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("sw-services-go listening on %s (context path %s)", addr, cfg.Server.ContextPath)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}


