package Routes

import (
	"seentech/RECR/Controllers"

	"github.com/gofiber/fiber/v2"
)

func CampaignRoute(route fiber.Router) {
	route.Post("/new", Controllers.CampaignNew)
	route.Put("/set_status/:campaignID/:status", Controllers.CampaignStatusModify)
	route.Put("/modify/:campaignID", Controllers.CampaignModify)
	route.Delete("/delete/:campaignID", Controllers.CampaignDelete)
	route.Post("/get_all", Controllers.CampaignGetAll)
	route.Post("/get_all_populated", Controllers.CampaignGetAllPopulated)
}
