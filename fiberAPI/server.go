package fiberapi

import (
	"baia_service/mongoService"
	"baia_service/utils"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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

	app.Get("/image/:id", func(c *fiber.Ctx) error {

		db := mongoClient.Database("Sushi_Restaurant")
		bucket, err := gridfs.NewBucket(db)
		if err != nil {
			panic(err)
		}

		idHex := c.Params("id")
		id, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid file ID")
		}

		downloadStream, err := bucket.OpenDownloadStream(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "File not found")
		}
		defer downloadStream.Close()

		c.Set("Content-Type", "image/jpeg") // Establece el tipo de contenido adecuado para la imagen
		_, err = io.Copy(c.Response().BodyWriter(), downloadStream)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error reading file")
		}

		return nil
	})

	return app
}
