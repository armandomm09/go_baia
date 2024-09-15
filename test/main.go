package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	Id          int    `json:"id"`
	ServiceName string `json:"service_name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Photo       string `json:"photo"`
}

type Menu struct {
	RestaurantName string    `json:"restaurant_name"`
	Services       []Service `json:"services"`
}

func main() {
	//************************************UNMARSHAL JSON
	jsonMenu, err := os.Open("menu.json")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")

	defer jsonMenu.Close()

	byteValue, _ := ioutil.ReadAll(jsonMenu)

	var menu Menu

	json.Unmarshal(byteValue, &menu)

	log.Println(menu.RestaurantName)

	// Configuración del cliente de MongoDB
	godotenv.Load()
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("No ENV variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// ************************************UPLOAD IMAGE
	db := client.Database("Sushi_Restaurant")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(menu.Services); i++ {
		imageURL := menu.Services[i].Photo

		// Descargar la imagen desde la URL
		resp, err := http.Get(imageURL)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Failed to download image: %s", resp.Status)
		}

		uploadOpts := options.GridFSUpload().SetMetadata(bson.D{{"metadata", "first"}})

		imageName := strings.ReplaceAll(menu.Services[i].ServiceName, " ", "_")
		objectID, err := bucket.UploadFromStream(fmt.Sprintf("%v.jpeg", imageName), resp.Body, uploadOpts)
		if err != nil {
			panic(err)
		}
		fmt.Printf("New file uploaded with ID %s\n", objectID)

		// Cambiar el JSON para usar la URL de la imagen en GridFS
		menu.Services[i].Photo = fmt.Sprintf("http://localhost:3000/image/%s", objectID.Hex())
	}

	// Guardar el JSON actualizado en un nuevo archivo
	updatedJSON, err := json.MarshalIndent(menu, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling updated JSON: %v", err)
	}

	if err := ioutil.WriteFile("updated_menu.json", updatedJSON, 0644); err != nil {
		log.Fatalf("Error writing updated JSON to file: %v", err)
	}

	fmt.Println("Updated JSON saved to updated_menu.json")

	// //************************************CREATE COLLECTION

	// // db := client.Database("BAIA")
	// // title := "Sushi Restaurant"
	// // command := bson.D{{"create", title}}

	// // var result bson.M

	// // if err := db.RunCommand(context.TODO(), command).Decode(&result); err != nil {
	// // 	log.Fatalln(err)
	// // }
	// // fmt.Println(time.Since(start))
	// // fmt.Printf("Collection %v created", title)

	// //************************************INSERT DOC

	// 	coll := client.Database("Sushi_Restaurant").Collection("Menu")

	// 	result, err := coll.InsertOne(context.TODO(), menu)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	// fmt.Println(time.Since(start))
	// 	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	// 	// //************************************FIND DOC
	// 	coll := client.Database("Sushi_Restaurant").Collection("Conversations")

	// 	// title := "Sakura Sushi"

	// 	var result bson.M

	// 	err = coll.FindOne(context.TODO(), bson.D{{"isActive", true}, {"userID", "5212223201384_c_us"}}).
	// 		Decode(&result)
	// 	if err == mongo.ErrNoDocuments {
	// 		fmt.Printf("No document was found")
	// 		return
	// 	}
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// jsonData, err := json.MarshalIndent(result, "", "    ")
	//
	//	if err != nil {
	//		panic(err)
	//	}
	//
	// fmt.Println(time.Since(start))
	// fmt.Printf("%s\n", jsonData)
}
