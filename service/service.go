package service

import (
	"context"

	"github.com/fmotalleb/pub-dev/config"
	"github.com/fmotalleb/pub-dev/web"
)

func Serve(ctx context.Context, cfg *config.Config) error {
	srv, err := web.NewServer(ctx, cfg)
	if err != nil {
		return err
	}
	return srv.Start()
}
