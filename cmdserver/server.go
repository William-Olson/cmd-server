package cmdserver

import (
	"github.com/labstack/echo"
	"github.com/william-olson/cmd-server/cmddeps"
	"github.com/william-olson/cmd-server/cmdutils"
)

// Server is the echo server manager
type Server struct {
	Deps *cmddeps.Deps
	e    *echo.Echo
}

type someRoutes interface {
	mount()
}

type route struct {
	group  *echo.Group
	deps   *cmddeps.Deps
	logger cmdutils.Logger
}

// Start inits routes and starts the server
func (s *Server) Start() {

	s.e = echo.New()
	sLogger := cmdutils.NewLogger("api")
	rLogger := cmdutils.NewLogger("api:routes")
	config := s.Deps.Get("config").(*cmdutils.Config)

	rootGroup := s.e.Group("")

	// define base paths
	routes := []someRoutes{
		rootRoutes{route{rootGroup, s.Deps, rLogger}},
	}

	// wire up sub paths
	for _, r := range routes {
		r.mount()
	}

	addr := ":"
	port := config.Get("APP_PORT")

	sLogger.
		KV("addr", addr).
		KV("port", port).
		Log("starting server")

	err := s.e.Start(addr + port)

	if err != nil {
		sLogger.
			KV("addr", addr).
			KV("port", port).
			KV("err", err).
			Error("server fatal error")
	}

}
