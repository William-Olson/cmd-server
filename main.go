package main

import (
	"github.com/william-olson/cmd-server/cmddeps"
	"github.com/william-olson/cmd-server/cmdserver"
	"github.com/william-olson/cmd-server/cmdutils"
)

func main() {

	deps := cmddeps.NewDeps()

	// init and set config
	config := cmdutils.NewConfig()
	deps.Set("config", &config)

	// start up the server
	server := cmdserver.Server{Deps: &deps}
	server.Start()

}
