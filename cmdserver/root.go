package cmdserver

import (
	"github.com/labstack/echo"
)

type rootRoutes struct {
	route
}

/*

	Map root routes

*/
func (r rootRoutes) mount() {

	r.group.GET("/", r.getRoot)

}

/*

	Serve the root route

*/
func (r rootRoutes) getRoot(c echo.Context) error {

	return c.JSON(200, map[string]bool{
		"ok": true,
	})

}
