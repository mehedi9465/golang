package Routes

import (
	"seentech/RECR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func TicketRoute(route fiber.Router) {
	route.Post("/new", Controllers.TicketNew)
	route.Put("/close/:ticketid/:reason/:deal", Controllers.TicketClose)
	route.Get("/get_all", Controllers.TicketsGetAll)
	route.Get("/get_all_populated", Controllers.TicketsGetAllPopulated)
}
