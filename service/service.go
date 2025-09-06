package service

import (
	"context"

	"github.com/fmotalleb/pub-dev/web"
)

func Serve(ctx context.Context) error {
	if err := web.Start(ctx); err != nil {
		return err
	}
	return nil
}
