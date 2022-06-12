package Routes

import (
	"seentech/RECR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func ProductRoute(route fiber.Router) {
	route.Post("/new", Controllers.ProductNew)
	route.Get("/get_all", Controllers.ProductsGetAll)
}
