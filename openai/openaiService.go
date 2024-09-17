package myOpenAi

import (
	"baia_service/mongoService"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// var jsonMenuData []byte
// var menuError error
// var jsonOrdersData []byte
// var ordersDataError error
// var fbClient *db.Client

// func InitOpenaiService(jsonMenuData []byte, jsonOrdersData []byte, newFbClient *db.Client) {
// 	jsonMenuData, menuError = ioutil.ReadFile("jsons/menu.json")
// 	if menuError != nil {
// 		fmt.Println("Error at parsing menu json")
// 	}

// 	jsonOrdersData, ordersDataError = ioutil.ReadFile("jsons/orders/order.json")
// 	if ordersDataError != nil {
// 		fmt.Println("Error at parsing order json")
// 	}

// 	fbClient = newFbClient
// }

func AskGpt(message string, userID string, mongoClient *mongo.Client) string {

	jsonMenuData, err := ioutil.ReadFile("jsons/menu.json")
	if err != nil {
		fmt.Println("Error at parsing menu json")
	}

	jsonOrdersData, err := ioutil.ReadFile("jsons/orders/order.json")
	if err != nil {
		fmt.Println("Error at parsing order json")
	}

	jsonOrderFotoData, err := ioutil.ReadFile("jsons/orders/fotoResultEG.json")
	if err != nil {
		fmt.Println("Error at parsing order as json")
	}

	Req := openai.ChatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `Eres un asistente de restaurante. Tu tarea es recibir pedidos, compararlos con el menú y devolver el pedido en formato JSON. 
				Cada platillo en la orden debe tener los campos: id, serviceName, unitaryPrice y quantity. Devuelve el JSON en el formato:` + string(jsonOrdersData) + `. 
				Si no hay pedido, la parte de "orden" regrasala vacia.
				Cuando se solicite el menú, no devuelvas el JSON del menú directamente, en un mensaje enlista todo el menu. Aquí está el menú:` + string(jsonMenuData) + `.
				Siempre responde en español. No resumas la orden en tus respuestas. Es decir, no hagas un resumen o un desglose de la orden 
				simplemente muestra el pedido o el menú según se solicite. Cuando el array de message no vaya vacio, haz un mensaje del total, y en
				ese mensaje, "afterOrder" seria true como en el ejemplo que te di.
				Recuerda que siempre tienes que hacer que el cliente compre lo mas posible, ofrecele diferentes cosas por ejemplo que combinen con lo que 
				te acaban de pedir.
				Cuando el usuario te diga que ya es todo o que ya completo su pedido, enseñale un resumen de la orden, el total, y literalmente di: "ORDEN COMPLETA".
				Si el usuario desea agregar algo más, muestra la orden actualizada y sugiéreles algo adicional para acompañar su pedido. 
				Cuando envíes fotos, solo incluye el enlace y establece isImage en true, sin texto adicional. Ejemplo:` + string(jsonOrderFotoData) + `
				Recuerda siempre mandar mensajes casuales que no siempre sea lo mismo, se muy emotivo y de vez en cuando usa emojis y siempre incita al cliente a comprar.`,
			},
		},
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// responseFormat := openai.ChatCompletionResponseFormat{
	// 	Type: openai.ChatCompletionResponseFormatTypeJSONObject,
	// }

	// Req.ResponseFormat = &responseFormat
	Req.ResponseFormat = &openai.ChatCompletionResponseFormat{
		Type: openai.ChatCompletionResponseFormatTypeJSONObject,
	}
	// Req = realtimeService.GetUserChatHistory(senderID, fbClient, Req)
	Req = mongoService.GetUserChatHistory(userID, mongoClient, Req)
	Req.Messages = append(Req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	resp, err := client.CreateChatCompletion(context.Background(), Req)
	if err != nil {
		fmt.Println("There was an error" + err.Error())
		return ""
	}
	log.Println("GPT RESPONSE:", resp.Choices[0].Message.Content)
	return resp.Choices[0].Message.Content
}

func Speech_to_text(filePathName string) string {
	openai_api_key := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(openai_api_key)
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filePathName,
		Language: "es",
	}

	// responseFormat := openai.ChatCompletionResponseFormat{
	// 	Type: openai.ChatCompletionResponseFormatTypeJSONObject,
	// }

	// Req.ResponseFormat = &responseFormat
	resp, err := client.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return ""
	}

	return string(resp.Text)
}
