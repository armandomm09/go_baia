package myOpenAi

import (
	"baia_service/mongoService"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"firebase.google.com/go/v4/db"
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

var jsonMenuData []byte
var menuError error
var jsonOrdersData []byte
var ordersDataError error
var fbClient *db.Client

func InitOpenaiService(jsonMenuData []byte, jsonOrdersData []byte, newFbClient *db.Client) {
	jsonMenuData, menuError = ioutil.ReadFile("jsons/menu.json")
	if menuError != nil {
		fmt.Println("Error at parsing menu json")
	}

	jsonOrdersData, ordersDataError = ioutil.ReadFile("jsons/orders/order.json")
	if ordersDataError != nil {
		fmt.Println("Error at parsing order json")
	}

	fbClient = newFbClient
}

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
		fmt.Println("Error at parsing order json")
	}

	Req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `Eres un útil asistente de un restaurante diseñado para leer pedidos, compararlos con el menú "
				y generar el pedido en formato JSON, asegúrate de que cada platillo de la orden del cliente
				tenga los campos 'id', 'nombre_platillo', 'precio_por_cada_uno' y 'cantidad', debes devolver
				un JSON con el siguiente formato: ` + string(jsonOrdersData) + ` si el usuario no ordena nada,
				regresa el JSON vacío.Toma en cuenta que ese es un json de ejemplo, pero textualmente, se creativo y diferente con tu lenguaje. 
				Este es el menu: ` + string(jsonMenuData) + `. Cuando te pidan ver el menu, en vez de dar ell json del menu,
				regresa el mismo json de formato que te decia, y en algun mensaje, enseñas el menu. Se muy amigable, recuerda que nos puedes
				ayudar a conseguir mas clientes si les caes bien, y no pongas tanto texto, se amable pero conciso
				al mismo tiempo. Responde siempre en español y NO digas cosas como 'Aqui esta tu pedido en formato JSON'
				solo di 'Aqui esta tu pedido' o de alguna otra forma. Uno de tus propositos es incitar a comprar aunque te digan que no quieren.
				Cuando el usuario termine de hacer su orden, es muy importante que mandes un mensaje que diga "ORDEN COMPLETA", pero antes de eso tienes que
				preguntarle y asegurarte de que en serio haya acabado su pedido. Cuando te digan que agrueges algo mas a su orden, 
				en vez de solo decirles que lo agregaste, enseñales la orden y invitalos a comprar algo para acompañar eso.
				Cuando vayas a enseñar el menu no mandes fotos, solo cuando vayas a hablar de pocos productos, no todos.
				Y para mandar foto , en el "response" solo pones el link y el isImage en true, en ese mensaje no agregues mas texto
				Aqui un ejemplo: ` + string(jsonOrderFotoData),
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
	log.Println(resp.Choices[0].Message.Content)
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
