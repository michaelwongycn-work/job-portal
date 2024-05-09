package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/michaelwongycn/job-portal/controller"
	"github.com/michaelwongycn/job-portal/handler"
	"github.com/michaelwongycn/job-portal/lib/auth"
	"github.com/michaelwongycn/job-portal/lib/cache"
	"github.com/michaelwongycn/job-portal/lib/cfg"
	"github.com/michaelwongycn/job-portal/lib/db"
	"github.com/michaelwongycn/job-portal/lib/encrypt"
	"github.com/michaelwongycn/job-portal/repository/appDB"
	"github.com/michaelwongycn/job-portal/usecase/job"
	"github.com/michaelwongycn/job-portal/usecase/user"
)

// @title job-portal
// @version 1.0.0
func main() {
	cfg, err := cfg.ReadConfig()
	if err != nil {
		log.Printf("Error reading config: %v\n", err)
	}

	auth.SetAuthConfig(cfg.JWT.SecretKey, cfg.JWT.AccessTokenDuration, cfg.JWT.RefreshTokenDuration)
	encrypt.SetAuthConfig(cfg.Encrypt.SecretKey)
	cache.InitializeNewCache(cfg.JWT.AccessTokenDuration)
	db, err := db.Connect(cfg.Database.Timeout, cfg.Database.DBName, cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password)
	if err != nil {
		log.Printf("Error connecting to DB: %v\n", err)
	}

	appDB := appDB.NewAppDBImpl(60, db)

	userUsecase := user.NewUserImpl(appDB, cfg.JWT.RefreshTokenDuration)
	JobUsecase := job.NewJobImpl(appDB, cfg.JWT.RefreshTokenDuration)

	controller := controller.NewControllerImpl(userUsecase, JobUsecase)

	handler := handler.NewHandler(60, controller)

	rest := handler.StartRoute()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Printf("Shutdown Application ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		cancel()
		db.Close()
	}()

	if err := rest.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown: %v", err)
	}
	log.Printf("Application Stopped")
}
