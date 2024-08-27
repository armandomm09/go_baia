package fiberapi

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	app.Get("/orders/:service/active", func(c *fiber.Ctx) error {
		service := c.Params("service") // Almacenar el valor antes de entrar en el StreamWriter

		c.Set("Content-Type", "text/event-stream; charset=utf-8")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")
		var currentJson []byte
		c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			fmt.Println("WRITER")

			for {
				// Usa la variable almacenada en lugar de acceder a c.Params directamente
				db := mongoClient.Database(service)
				coll := db.Collection("Orders")

				var results []bson.M

				cursor, err := coll.Find(
					context.TODO(),
					bson.M{"state": "active"},
					options.Find().SetSort(bson.D{{"creationDate", -1}}),
				)
				if err != nil {
					fmt.Println("Error al buscar las órdenes:", err)
					break
				}
				defer cursor.Close(context.TODO())

				if err = cursor.All(context.TODO(), &results); err != nil {
					fmt.Println("Error al decodificar las órdenes:", err)
					break
				}

				jsonData, err := json.MarshalIndent(results, "", "    ")
				if err != nil {
					fmt.Println("Error al convertir a JSON:", err)
					break
				}
				if !bytes.Equal(currentJson, jsonData) {
					fmt.Println("DIF")
					if _, err := w.Write([]byte("" + string(jsonData) + "\n\n")); err != nil {
						fmt.Println("Error al escribir los datos:", err)
						break
					}
					currentJson = jsonData
					// fmt.Println("Datos enviados:", string(jsonData))
				}
				if err := w.Flush(); err != nil {
					fmt.Printf("Error al hacer flush: %v. Cerrando la conexión HTTP.\n", err)
					break
				}

				// time.Sleep(2 * time.Microsecond)
			}
		}))
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

	//Gets the specific image using ID
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
