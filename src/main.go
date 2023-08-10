package main

import (
	"vnc-write-api/api/config"
	_ "vnc-write-api/docs"
)

// @Title       VNC Write API
// @Version     v1
// @Description Este repositório é responsável pela escrita dos dados nas bases de dados da Plataforma Você na Câmara.
// @BasePath    /api/v1
func main() {
	config.NewServer()
}
