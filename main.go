package main

import (
	"github.com/william-olson/cmd-server/cmddb"
	"github.com/william-olson/cmd-server/cmddeps"
	"github.com/william-olson/cmd-server/cmdserver"
	"github.com/william-olson/cmd-server/cmdutils"
)

func main() {

	deps := cmddeps.NewDeps()
	logger := cmdutils.NewLogger("app")
	depsSet := []string{}

	// init and set config
	config := cmdutils.NewConfig()
	deps.Set("config", &config)
	depsSet = append(depsSet, "config")

	// init and set db
	db := cmddb.DB{Deps: &deps}
	db.Connect()
	deps.Set("db", &db)
	depsSet = append(depsSet, "db")

	// log set dependencies
	logger.KV("deps", depsSet).Log("registered dependencies")

	// start up the server
	server := cmdserver.Server{Deps: &deps}
	server.Start()

}
