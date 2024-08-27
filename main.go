package main

import (
	fiberapi "baia_service/fiberAPI"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// 	"baia_service/utils"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Options struct {
	port int `help:"Port to listen on" short:"p" default:"8888"`
}

type GPTResponse struct {
	Body struct {
		Answer string `json:answer`
	}
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Obtener la dirección IP del cliente
		clientIP := r.RemoteAddr
		if ip := r.Header.Get("X-Real-IP"); ip != "" {
			clientIP = ip
		} else if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
			clientIP = ip
		}

		// Imprimir detalles de la solicitud
		fmt.Println("\n- - - - - - - INCOMING REQUEST - - - - - - - -\n")
		fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
		fmt.Printf("Client IP: %s\n", clientIP)
		fmt.Printf("User Agent: %s\n", r.UserAgent())
		// fmt.Printf("Headers:\n")
		// for name, values := range r.Header {
		// 	for _, value := range values {
		// 		fmt.Printf("  %s: %s\n", name, value)
		// 	}
		// }

		next.ServeHTTP(w, r)

		fmt.Printf("Completed request: %s %s in %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

type GPTRequest struct {
	Body struct {
		Question string `json:question`
	}
}

func main() {
	godotenv.Load()
	//************MONGO

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("No ENV variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Permitir todos los orígenes (no recomendado para producción)
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,OPTIONS",
		AllowHeaders: "*", //Origin, Content-Type, Accept
	}))

	app = fiberapi.RegisterEndPoints(app, client)

	app.Listen("10.50.94.111:8000")

}

//***********************FIREBASE AND HUMA
// func main() {
// 	godotenv.Load()

// 	rtClient := realtimeService.InitFirebase()
// 	fbClient := firestoreService.InitFirebase()
// 	jsonMenuData, err := ioutil.ReadFile("jsons/menu.json")
// 	if err != nil {
// 		fmt.Println("Error at parsing menu json")
// 	}

// 	jsonOrdersData, err := ioutil.ReadFile("jsons/orders/order.json")
// 	if err != nil {
// 		fmt.Println("Error at parsing order json")
// 	}

// 	myOpenAi.InitOpenaiService(jsonMenuData, jsonOrdersData, rtClient)
// 	if err != nil {
// 		fmt.Println("Error initializing Firebase")
// 	}

// 	cli := humacli.New(func(hook humacli.Hooks, options *Options) {
// 		router := chi.NewMux()
// 		router.Use(requestLogger)
// 		api := humachi.New(router, huma.DefaultConfig("My First API", "1.0.0"))

// 		hook.OnStart(func() {
// 			fmt.Printf("Starting server on port %d...\n", 8888)
// 			http.ListenAndServe(fmt.Sprintf(":%d", 8888), router)
// 		})
// 		baiaAPI.RegisterEndPoints(api, rtClient, fbClient)

// 	})

// 	cli.Run()
// }
