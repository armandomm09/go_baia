package main

import (
	baiaAPI "baia_service/api"
	"baia_service/firebase/firestoreService"
	"baia_service/firebase/realtimeService"
	myOpenAi "baia_service/openai"
	"fmt"
	"io/ioutil"
	"time"

	// 	"baia_service/utils"

	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
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
	rtClient := realtimeService.InitFirebase()
	fbClient := firestoreService.InitFirebase()
	jsonMenuData, err := ioutil.ReadFile("jsons/menu.json")
	if err != nil {
		fmt.Println("Error at parsing menu json")
	}

	jsonOrdersData, err := ioutil.ReadFile("jsons/orders/order.json")
	if err != nil {
		fmt.Println("Error at parsing order json")
	}

	myOpenAi.InitOpenaiService(jsonMenuData, jsonOrdersData, rtClient)
	if err != nil {
		fmt.Println("Error initializing Firebase")
	}

	cli := humacli.New(func(hook humacli.Hooks, options *Options) {
		router := chi.NewMux()
		router.Use(requestLogger)
		api := humachi.New(router, huma.DefaultConfig("My First API", "1.0.0"))

		hook.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", 8888)
			http.ListenAndServe(fmt.Sprintf(":%d", 8888), router)
		})
		baiaAPI.RegisterEndPoints(api, rtClient, fbClient)

	})

	cli.Run()
}
