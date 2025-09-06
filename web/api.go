package web

import (
	"github.com/fmotalleb/pub-dev/web/api"
)

func init() {
	RegisterEndpoint(api.Setup)
}
