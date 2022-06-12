package Routes

import (
	"seentech/RECR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func ChannelRoute(route fiber.Router) {
	route.Post("/new", Controllers.ChannelNew)
	route.Get("/get_all/", Controllers.ChannelGet)
	route.Put("/set_status/:id/:new_status", Controllers.ChannelSetStatus)
	route.Put("/modify", Controllers.ChannelModify)
}
