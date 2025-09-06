package service

import (
	"context"

	"github.com/fmotalleb/go-tools/log"
	"go.uber.org/zap"

	"github.com/fmotalleb/pub-dev/config"
	"github.com/fmotalleb/pub-dev/database"
	"github.com/fmotalleb/pub-dev/models"
	"github.com/fmotalleb/pub-dev/web"
)

func Serve(ctx context.Context) error {
	l := log.FromContext(ctx).Named("Serve")
	cfg, err := config.Get(ctx)
	if err != nil {
		return err
	}
	db, err := database.Connect(cfg.DatabaseConnection)
	if err != nil {
		l.Error("failed to connect to db", zap.Error(err))
		return err
	}
	if err = db.AutoMigrate(models.GetModels()...); err != nil {
		l.Error("failed to migrate db", zap.Error(err))
		return err
	}

	if err := web.Start(ctx); err != nil {
		return err
	}
	return nil
}
