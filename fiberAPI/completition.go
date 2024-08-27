package fiberapi

import (
	"baia_service/mongoService"
	"baia_service/utils"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type SendMessageInputT struct {
	Question string `json:"question" example:"Hola"`
	SenderID string `json:"senderID" example:"5212223201384@c.us"`
}

type OutputMessage struct {
	Response string `json:"response"`
	IsImage  bool   `json:"isImage"`
}

type SendMessageOutputT struct {
	Messages []OutputMessage `json:"messages"`
}

func RegisterEndPoints(app *fiber.App, mongoClient *mongo.Client) *fiber.App {
	app = RegisterDBEndPoints(app, mongoClient)

	app.Post("/baia/askGPT/text", func(c *fiber.Ctx) error {
		start := time.Now()
		var input SendMessageInputT
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}
		answer, actualOrder := utils.SendRequest(input.Question, input.SenderID, mongoClient)
		if strings.Contains(fmt.Sprintf("%v", answer), "ORDEN COMPLETA") {
			mongoService.FinishOrder("Sushi_Restaurant", input.SenderID, actualOrder, mongoClient)
			log.Println(" ORDEN COMPLETA")
		}
		log.Printf("Tiempo de respuesta: %v", time.Since(start))
		return c.JSON(answer)
	})

	return app
}
