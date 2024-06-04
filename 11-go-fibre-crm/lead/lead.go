package lead

import (
	"crm/database"

	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Lead struct {
	gorm.Model
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   int    `json:"phone"`
}

func GetLeads(context *fiber.Ctx) {
	db := database.Connection
	var leads []Lead
	db.Find(&leads)
	context.JSON(leads)
}

func GetLead(context *fiber.Ctx) {
	id := context.Params("id")
	db := database.Connection
	var lead Lead
	db.Find(&lead, id)
	context.JSON(lead)
}

func NewLead(context *fiber.Ctx) {
	db := database.Connection
	lead := new(Lead)
	if err := context.BodyParser(lead); err != nil {
		context.Status(503).Send(err)
		return
	}
	db.Create(&lead)
	context.JSON(lead)
}

func DeleteLead(context *fiber.Ctx) {
	id := context.Params("id")
	db := database.Connection
	var lead Lead
	db.First(&lead, id)
	if lead.Name == "" {
		context.Status(500).Send("No lead found.")
		return
	}
	db.Delete(&lead)
	context.Send("Deleted lead successfully")
}
