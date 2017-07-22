package cmdserver

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/william-olson/cmd-server/cmddb"
	"github.com/william-olson/cmd-server/cmdversions"
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
	r.group.GET("/version", r.getVersion)
	r.group.POST("/slug-version", r.getSlugs)

}

/*

	Serve the root route

*/
func (r rootRoutes) getRoot(c echo.Context) error {

	return c.JSON(200, map[string]bool{
		"ok": true,
	})

}

/*

	Test route to view slack_clients

*/
func (r rootRoutes) getClients(c echo.Context) error {

	db := r.deps.Get("db").(*cmddb.DB)
	clients, err := db.GetSlackClientsOrErr()

	if err != nil {
		fmt.Printf("%v\n", err)
		return c.String(500, "Error")
	}

	return c.JSON(200, clients)

}

/*

	Fetch the version for the default server

*/
func (r rootRoutes) getVersion(c echo.Context) error {

	resp, err := cmdversions.GetDefaultOrErr(r.deps)

	if err != nil {
		fmt.Println(err)
		return c.String(500, "Error")
	}

	return c.JSON(200, resp)

}

/*

	Fetch a slug version based on slack_client

	Restrictions:
		- slack_client with token must be in db
		- must have version_path and host fields set

*/
func (r rootRoutes) getSlugs(c echo.Context) error {

	slug := "all"
	db := r.deps.Get("db").(*cmddb.DB)

	// extract the token and text args
	token := c.FormValue("token")
	text := c.FormValue("text")

	if len(text) != 0 {
		slug = text
	}

	// fetch the slack_client from db
	slackClient, err := db.GetSlackClientByTokenOrErr(token)

	if err != nil {
		c.JSON(400, map[string]string{"error": "Bad Token"})
	}

	// TODO:
	//  check if slash command arguments contain multiple slugs

	resp, slugErr := cmdversions.GetSlugVersionOrErr(db, slackClient, slug)

	if slugErr != nil {
		c.JSON(500, map[string]string{"error": "Version fetching problem"})
	}

	return c.JSON(200, resp)

}
