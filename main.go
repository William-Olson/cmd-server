package main

import (
	"github.com/william-olson/cmd-server/cmddb"
	"github.com/william-olson/cmd-server/cmddeps"
	"github.com/william-olson/cmd-server/cmdserver"
	"github.com/william-olson/cmd-server/cmdutils"
)

func main() {

	deps := cmddeps.NewDeps()
	config := cmdutils.NewConfig()

	// setup logger
	logToken := config.Get("SENTRY_TK")
	logger := cmdutils.NewLogger("app")
	if len(logToken) > 0 {
		logProj := config.Get("SENTRY_PRJ")
		logger.EnableAggregation(logToken + " " + logProj)
	}

	// init and set config
	logger.KV("dep", "config").Log("registering depedency")
	deps.Set("config", &config)

	// init and set db
	logger.KV("dep", "db").Log("registering depedency")
	db := cmddb.DB{Deps: &deps}
	db.Connect()
	deps.Set("db", &db)

	// start up the server
	server := cmdserver.Server{Deps: &deps}
	server.Start()

}
