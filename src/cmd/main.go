package main

import (
	"kudejen/src/internal/api/interceptor"
	"kudejen/src/internal/api/server"
)

func main() {
	srv := server.NewServer("src/internal/config/config.yml", "src/internal/config/kube.config")
	interceptor.AddMiddlewares(srv)
}
