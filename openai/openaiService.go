package myOpenAi

import (
	"baia_service/firebase/realtimeService"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"firebase.google.com/go/v4/db"
	"github.com/sashabaranov/go-openai"
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

func AskGpt(message string, senderID string) string {

	jsonMenuData, err := ioutil.ReadFile("jsons/menu.json")
	if err != nil {
		fmt.Println("Error at parsing menu json")
	}

	jsonOrdersData, err := ioutil.ReadFile("jsons/orders/order.json")
	if err != nil {
		fmt.Println("Error at parsing order json")
	}

	exampleResponse := "¡Perfecto! Entonces tenemos:- **1 Miso Soup**- **2 Sakes**Aquí está tu pedido:```json{    " + "orden" + ": [        {            " + "id" + ": 8,            " + "nombre_platillo" + ": " + "Miso Soup" + ",            " + "precio_por_cada_uno" + ": 5.00,            " + "cantidad" + ": 1        },        {            " + "id" + ": 12,            " + "nombre_platillo" + ": " + "Sake" + ",            " + "precio_por_cada_uno" + ": 9.00,            " + "cantidad" + ": 2        }    ]}```¿Te gustaría agregar algo más para acompañar? ¡Aprovecha que tenemos delicias japonesas en nuestro menú!"
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
				Este es el menu: ` + string(jsonMenuData) + `Se muy amigable, recuerda que nos puedes
				ayudar a conseguir mas clientes si les caes bien, y no pongas tanto texto, se amable pero conciso
				al mismo tiempo. Responde siempre en español y NO digas cosas como 'Aqui esta tu pedido en formato JSON'
				solo di 'Aqui esta tu pedido' o de alguna otra forma. Uno de tus propositos es incitar a comprar aunque te digan que no quieren.
				Cuando el usuario termine de hacer su orden, es muy importante que mandes un mensaje que diga "ORDEN COMPLETA", pero antes de eso tienes que
				preguntarle y asegurarte de que en serio haya acabado su pedido un ejemplo de como terminas tu orden es asi: ` + exampleResponse,
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
	Req = realtimeService.GetUserChatHistory(senderID, fbClient, Req)

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
