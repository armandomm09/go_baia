package myOpenAi

import (
	firebaseService "baia_service/firebase"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"cloud.google.com/go/firestore"
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
var fbClient *firestore.Client

func InitOpenaiService(jsonMenuData []byte, jsonOrdersData []byte, newFbClient *firestore.Client) {
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

	Req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo1106,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `Eres un útil asistente de un restaurante diseñado para leer pedidos, compararlos con el menú "
				y generar el pedido en formato JSON, asegúrate de que cada platillo de la orden del cliente
				tenga los campos 'id', 'nombre_platillo', 'precio_por_cada_uno' y 'cantidad', debes devolver
				un JSON con el siguiente formato: ` + string(jsonOrdersData) + ` si el usuario no ordena nada,
				regresa el JSON vacío. Menu: ` + string(jsonMenuData) + `Se muy amigable, recuerda que nos puedes
				ayudar a conseguir mas clientes si les caes bien, y no pongas tanto texto, se amable pero conciso
				al mismo tiempo. Responde siempre en español y NO digas cosas como 'Aqui esta tu pedido en formato JSON'
				solo di 'Aqui esta tu pedido' o de alguna otra forma. Uno de tus propositos es incitar a comprar aunque te digan que no quieren`,
			},
		},
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	Req = firebaseService.GetUserChatHistory(senderID, fbClient, Req, jsonMenuData, jsonOrdersData)

	Req.Messages = append(Req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	resp, err := client.CreateChatCompletion(context.Background(), Req)
	if err != nil {
		fmt.Println("There was an error" + err.Error())
		return ""
	}

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
