package main

import (
	"context"
	"github.com/arkdbp/jobservice/server"
)

func main() {
	app := WireApp(server.ProvideConfig())
	app.Run(context.Background())
}
