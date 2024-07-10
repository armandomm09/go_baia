package firebaseService

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/sashabaranov/go-openai"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func GetUserChatHistory(userID string, client *firestore.Client, Req openai.ChatCompletionRequest, jsonMenuData []byte, jsonOrdersData []byte) openai.ChatCompletionRequest {
	ctx := context.Background()

	// Obtener todos los documentos de la colección "Conversation" para el usuario dado
	docs, err := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(userID).Doc("messages").Collection("Conversation").OrderBy("timestamp", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		fmt.Println("Error getting documents:", err)
		return Req
	}
	if len(docs) == 0 {
		SetInitialPromt(userID, jsonMenuData, jsonOrdersData, client)
		fmt.Println("********** INITIAL PROMT **********")
		return Req
	}
	// Verificar el número de documentos obtenidos
	// fmt.Printf("Number of documents retrieved: %d\n", len(docs))

	// Iterar sobre los documentos y agregarlos a Req.Messages
	for i, doc := range docs {
		var msg Message
		if err := doc.DataTo(&msg); err != nil {
			fmt.Printf("Error converting document data for document %d: %v\n", i, err)
			continue
		}

		// fmt.Printf("Appending message %d: Role=%s, Content=%s\n", i, msg.Role, msg.Content)

		Req.Messages = append(Req.Messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return Req
}
