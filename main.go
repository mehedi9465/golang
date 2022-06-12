package main

import (
	"fmt"
	"seentech/RECR/Controllers"
	"seentech/RECR/DBManager"
	"seentech/RECR/Routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func SetupRoutes(app *fiber.App) {
	Routes.UserRoute(app.Group("/user"))
	Routes.ProductRoute(app.Group("/product"))
	Routes.TicketRoute(app.Group("/ticket"))
	Routes.ChannelRoute(app.Group("/channel"))
	Routes.CampaignRoute(app.Group("/campaign"))
	Routes.SettingRoute(app.Group("/setting"))
}

func main() {
	fmt.Println("Hello RECR")
	fmt.Print("Initializing Database Connections ... ")
	initState := DBManager.InitCMSCollections()
	initSetting := Controllers.InitializeSetting()
	fmt.Println("Hello")

	if initState && initSetting {
		fmt.Println("[OK]")
	} else {
		fmt.Println("[FAILED]")
		return
	}

	fmt.Print("Initializing the server ... ")

	app := fiber.New()
	app.Use(cors.New())
	app.Use(pprof.New())
	SetupRoutes(app)
	app.Static("/Public", "./Public")

	fmt.Println("[OK]")
	app.Listen(":8081")

}
