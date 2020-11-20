package main

import (
	"fmt"
	"net/http"

	dbConn "api/adapter/gorm"
	"api/app/router"
	"api/app/server"
	"api/config"
	auth "api/util/auth"
	lr "api/util/logger"
	vr "api/util/validator"
)

func main() {
	appConf := config.AppConfig()

	logger := lr.New(appConf.Debug)

	db, err := dbConn.New(appConf)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
		return
	}
	if appConf.Debug {
		db.LogMode(true)
	}
	defer db.Close()

	validator := vr.New()

	jwtUtils := auth.New(appConf)

	server := server.New(logger, db, validator, jwtUtils)

	appRouter := router.New(server, appConf)

	address := fmt.Sprintf(":%d", appConf.Server.Port)

	logger.Info().Msgf("Starting server %v", address)

	s := &http.Server{
		Addr:         address,
		Handler:      appRouter,
		ReadTimeout:  appConf.Server.TimeoutRead,
		WriteTimeout: appConf.Server.TimeoutWrite,
		IdleTimeout:  appConf.Server.TimeoutIdle,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server startup failed")
	}
}
