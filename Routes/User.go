package Routes

import (
	"seentech/RECR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(route fiber.Router) {
	route.Post("/new", Controllers.UserNew)
	route.Post("/get_all", Controllers.UserGetAll)
}
