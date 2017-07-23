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
	group *echo.Group
	deps  *cmddeps.Deps
}

// Start inits routes and starts the server
func (s *Server) Start() {

	s.e = echo.New()
	config := s.Deps.Get("config").(*cmdutils.Config)

	rootGroup := s.e.Group("")

	// define base paths
	routes := []someRoutes{
		rootRoutes{route{rootGroup, s.Deps}},
	}

	// wire up sub paths
	for _, r := range routes {
		r.mount()
	}

	s.e.Logger.Fatal(s.e.Start(":" + config.Get("APP_PORT")))

}
