package cmdserver

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/william-olson/cmd-server/cmddb"
)

type rootRoutes struct {
	route
}

/*

	Map root routes

*/
func (r rootRoutes) mount() {

	r.group.GET("/", r.getRoot)
	r.group.GET("/clients", r.getClients)

}

/*

	Serve the root route

*/
func (r rootRoutes) getRoot(c echo.Context) error {

	return c.JSON(200, map[string]bool{
		"ok": true,
	})

}

func (r rootRoutes) getClients(c echo.Context) error {

	db := r.deps.Get("db").(*cmddb.DB)
	clients, err := db.GetSlackClients()

	if err != nil {
		fmt.Printf("%v\n", err)
		return c.String(500, "Error")
	}

	return c.JSON(200, clients)

}
