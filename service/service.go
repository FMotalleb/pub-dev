package service

import (
	"context"

	"github.com/golang-templates/seed/web"
)

func Serve(ctx context.Context) error {
	web.Start(ctx)
	return nil
}
