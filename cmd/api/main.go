// @title VKarmane API
// @version 1.0
// @description API для управления личными финансами

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name auth_token
// @description JWT токен в HTTP-only cookie

package main

import (
	"github.com/go-park-mail-ru/2025_2_VKarmane/cmd/api/app"
	_ "github.com/go-park-mail-ru/2025_2_VKarmane/docs"
)

func main() {
	app.Run()
}
