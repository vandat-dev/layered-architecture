package main

import (
	_ "app/docs" // Import for swag to recognize json.RawMessage
	"app/internal/initialize"
	_ "encoding/json"
)

// @title Go API
// @version 1.0
// @description This is a sample server Go API server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT Authorization header using Bearer scheme. Example: "Bearer {token}"
func main() {
	initialize.Run()
}
