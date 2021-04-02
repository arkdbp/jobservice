//+build wireinject

package main

import (
	"github.com/arkdbp/jobservice/lib"
	"github.com/arkdbp/jobservice/server"
	"github.com/google/wire"
)

// WireApp for initialize app
func WireApp(c *server.Config) *server.App {
	wire.Build(server.ServiceSet, lib.CMDLibSet)

	return nil
}
