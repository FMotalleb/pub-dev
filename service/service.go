package service

import (
	"context"

	"github.com/fmotalleb/pub-dev/web"
)

func Serve(ctx context.Context) error {
	// l := log.FromContext(ctx).Named("Serve")
	// cfg, err := config.Get(ctx)
	// if err != nil {
	// 	return err
	// }

	if err := web.Start(ctx); err != nil {
		return err
	}
	return nil
}
