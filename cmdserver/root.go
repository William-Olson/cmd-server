package cmdserver

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/william-olson/cmd-server/cmddb"
	"github.com/william-olson/cmd-server/cmdutils"
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

	Fetch the version for the default server or a query url

*/
func (r rootRoutes) getVersion(c echo.Context) error {

	url := c.QueryParam("q")

	if len(url) > 0 {
		resp, err := cmdversions.GetVersionByUrlOrErr(url)
		if err != nil {
			return c.String(500, "Error")
		}
		return c.JSON(200, resp)
	}

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

	type CommandReq struct {
		Token string `json:"token"`
		Text  string `json:"text"`
		Cmd   string `json:"command"`
	}
	body := CommandReq{}

	bodyErr := c.Bind(&body)
	if bodyErr != nil {
		return c.JSON(400, map[string](interface{}){"error": bodyErr})
	}

	db := r.deps.Get("db").(*cmddb.DB)

	// fetch the slack_client from db
	slackClient, err := db.GetSlackClientByTokenOrErr(body.Token)

	if err != nil {
		return c.JSON(400, map[string]string{"error": "Bad Token"})
	}

	// check for multiple slug arguments
	slugs := cmdutils.SplitBySpaces(body.Text)

	// handle empty and all case
	if len(slugs) == 0 || slugs[0] == "all" || slugs[0] == "" {
		slugs = []string{}
		for _, slg := range slackClient.GetSlugs() {
			slugs = append(slugs, slg.Name)
		}
	}

	// ensure at least 1 will be requested
	if len(slugs) == 0 {
		return c.JSON(400, map[string]string{"error": "No slugs found"})
	}

	// handle single slug
	if len(slugs) == 1 {
		resp, slugErr := cmdversions.GetSlugVersionOrErr(db, slackClient, slugs[0])

		if slugErr != nil {
			return c.JSON(500, map[string]string{"error": "Version fetching problem"})
		}

		return c.JSON(200, resp)
	}

	// otherwise handle the multi-slugs
	resp, multiErr := cmdversions.GetMultiSlugVersionsOrErr(db, slackClient, slugs)

	if multiErr != nil {
		fmt.Println(multiErr)
		return c.JSON(500, map[string]string{"error": "Version fetching problem"})
	}

	return c.JSON(200, resp)

}
