package api

import "github.com/fmotalleb/pub-dev/web/api/packages"

func init() {
	RegisterEndpoint(packages.Setup)
}
