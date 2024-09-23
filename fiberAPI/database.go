package fiberapi

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RegisterDBEndPoints(app *fiber.App, mongoClient *mongo.Client) *fiber.App {

	// @Summary Obtiene todas las órdenes activas
	// @Description Retorna una lista de órdenes con estado "active" para un servicio específico
	// @Tags Orders
	// @Param service path string true "Nombre del servicio"
	// @Produce json
	// @Success 200 {array} bson.M
	// @Failure 500 {string} string "Error al buscar o decodificar las órdenes"
	// @Router /{service}/orders/active [get]
	app.Get("/:service/orders/active", func(c *fiber.Ctx) error {
		service := c.Params("service")

		// Conectar a la base de datos
		db := mongoClient.Database(service)
		coll := db.Collection("Orders")

		// Buscar las órdenes con estado "active"
		var results []bson.M
		cursor, err := coll.Find(
			context.TODO(),
			bson.M{"state": "active"},
			options.Find().SetSort(bson.D{{"creationDate", -1}}),
		)
		if err != nil {
			fmt.Println("Error al buscar las órdenes:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error al buscar las órdenes")
		}

		// Convertir los resultados en una lista
		if err = cursor.All(context.TODO(), &results); err != nil {
			fmt.Println("Error al decodificar las órdenes:", err)
			cursor.Close(context.TODO())
			return c.Status(fiber.StatusInternalServerError).SendString("Error al decodificar las órdenes")
		}
		cursor.Close(context.TODO()) // Cerrar el cursor después de la consulta

		// Convertir los resultados a JSON
		jsonData, err := json.MarshalIndent(results, "", "    ")
		if err != nil {
			fmt.Println("Error al convertir a JSON:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error al convertir a JSON")
		}

		// Devolver los datos en la respuesta
		c.Set("Content-Type", "application/json")
		return c.Status(fiber.StatusOK).Send(jsonData)
	})

	// @Summary Finaliza una orden
	// @Description Cambia el estado de una orden a "inactive" por ID en un servicio específico
	// @Tags Orders
	// @Param service path string true "Nombre del servicio"
	// @Param id path string true "ID de la orden"
	// @Produce text/plain
	// @Success 200 {string} string "Order Finished"
	// @Failure 500 {string} string "Error at updating doc"
	// @Router /{service}/orders/{id}/finishOrder [patch]
	app.Patch("/:service/orders/:id/finishOrder", func(c *fiber.Ctx) error {
		coll := mongoClient.Database("Sushi_Restaurant").Collection("Orders")
		oid := c.Params("id") // Almacenar el valor antes de entrar en el StreamWriter
		log.Println(oid)
		filter := bson.D{{"ID", oid}}

		update := bson.D{{"$set", bson.D{{"state", "inactive"}}}}

		result, err := coll.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return c.SendString(fmt.Sprintf("Error at updating doc: %v", err))
			panic(err)
		}
		log.Println(result)
		return c.SendString("Order Finished")
	})

	// @Summary Reactiva una orden
	// @Description Cambia el estado de una orden a "active" por ID en un servicio específico
	// @Tags Orders
	// @Param service path string true "Nombre del servicio"
	// @Param id path string true "ID de la orden"
	// @Produce text/plain
	// @Success 200 {string} string "Order Reactivated"
	// @Failure 500 {string} string "Error at updating doc"
	// @Router /{service}/orders/{id}/reactivateOrder [patch]
	app.Patch("/:service/orders/:id/reactivateOrder", func(c *fiber.Ctx) error {
		coll := mongoClient.Database("Sushi_Restaurant").Collection("Orders")
		oid := c.Params("id") // Almacenar el valor antes de entrar en el StreamWriter
		log.Println(oid)
		filter := bson.D{{"ID", oid}}

		update := bson.D{{"$set", bson.D{{"state", "active"}}}}

		result, err := coll.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			return c.SendString(fmt.Sprintf("Error at updating doc: %v", err))
			panic(err)
		}
		log.Println(result)
		return c.SendString("Order Reactivated")
	})

	// @Summary Obtiene una imagen por ID
	// @Description Descarga una imagen de GridFS usando su ID
	// @Tags Images
	// @Param id path string true "ID de la imagen"
	// @Produce image/jpeg
	// @Success 200 {file} file
	// @Failure 400 {string} string "Invalid file ID"
	// @Failure 404 {string} string "File not found"
	// @Failure 500 {string} string "Error reading file"
	// @Router /image/{id} [get]
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

	app.Get("/sse", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			fmt.Println("WRITER")
			var i int
			for {
				i++
				msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
				fmt.Fprintf(w, "data: Message: %s\n\n", msg)
				fmt.Println(msg)

				err := w.Flush()
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)

					break
				}
				time.Sleep(2 * time.Second)
			}
		}))

		return nil
	})
	return app
}
