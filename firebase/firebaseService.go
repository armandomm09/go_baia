package firebaseService

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
)

func InitFirebase() (*firestore.Client, error) {
	ctx := context.Background()
	config := &firebase.Config{
		ProjectID: "baia-1df5a",
	}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return client, nil
}

func getNextMessageNumber(client *firestore.Client, senderID string, isFromUser bool, isRaw bool) int {
	ctx := context.Background()

	var messageType string
	if isFromUser {
		messageType = "message"
	} else {
		messageType = "response"
	}

	var isRawMessage string
	if isRaw {
		isRawMessage = "rawConversation"
	} else {
		isRawMessage = "Conversation"
	}
	collectionRef := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc("messages").Collection(isRawMessage)
	iter := collectionRef.Documents(ctx)
	messageCount := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		docID := doc.Ref.ID
		if strings.Contains(docID, messageType) {
			messageCount++

		}
	}
	// fmt.Printf("%v number %v", messageType, messageCount+1)
	return messageCount + 1
}

func SaveRawBAIAMessage(message string, senderID string, client *firestore.Client) {
	ctx := context.Background()
	messageNumber := getNextMessageNumber(client, senderID, false, true)
	ref1 := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc("messages")
	ref := ref1.Collection("rawConversation").Doc(fmt.Sprintf("#%v response", messageNumber))
	_, err := ref.Set(ctx, map[string]interface{}{
		"content":   message,
		"role":      "assistant",
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		log.Fatalln(err)
	}

	// fmt.Println(result)
}

func SaveRawUserMessage(message string, senderID string, client *firestore.Client) {
	ctx := context.Background()
	messageNumber := getNextMessageNumber(client, senderID, true, true)
	ref1 := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc("messages")
	ref := ref1.Collection("rawConversation").Doc(fmt.Sprintf("#%v message", messageNumber))
	_, err := ref.Set(ctx, map[string]interface{}{
		"content":   message,
		"role":      "user",
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		log.Fatalln(err)
	}

	// fmt.Println(result)
}

func SaveBAIAMessage(message string, senderID string, client *firestore.Client) {
	ctx := context.Background()
	messageNumber := getNextMessageNumber(client, senderID, false, false)
	ref1 := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc("messages")
	ref := ref1.Collection("Conversation").Doc(fmt.Sprintf("#%v response", messageNumber))
	_, err := ref.Set(ctx, map[string]interface{}{
		"content":   message,
		"role":      "assistant",
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		log.Fatalln(err)
	}

	// fmt.Println(result)
}

func SaveUserMessage(message string, senderID string, client *firestore.Client) {
	ctx := context.Background()

	messageNumber := getNextMessageNumber(client, senderID, true, false)
	ref1 := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc("messages")
	ref := ref1.Collection("Conversation").Doc(fmt.Sprintf("#%v message", messageNumber))
	_, err := ref.Set(ctx, map[string]interface{}{
		"content":   message,
		"role":      "user",
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		log.Fatalln(err)
	}

}

func SetInitialPromt(senderID string, jsonMenuData []byte, jsonOrdersData []byte, client *firestore.Client) {
	ctx := context.Background()

	ref1 := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc("messages")
	ref := ref1.Collection("Conversation").Doc("Context promt")
	_, err := ref.Set(ctx, map[string]interface{}{
		"content": `Eres un útil asistente de un restaurante diseñado para leer pedidos, compararlos con el menú "
		y generar el pedido en formato JSON, asegúrate de que cada platillo de la orden del cliente
		tenga los campos 'id', 'nombre_platillo', 'precio_por_cada_uno' y 'cantidad', debes devolver
		un JSON con el siguiente formato: ` + string(jsonOrdersData) + ` si el usuario no ordena nada,
		regresa el JSON vacío. Menu: ` + string(jsonMenuData) + `Se muy amigable, recuerda que nos puedes
		ayudar a conseguir mas clientes si les caes bien, y no pongas tanto texto, se amable pero conciso
		al mismo tiempo. Responde siempre en español`,
		"role":      "system",
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		log.Fatalln(err)
	}

	ref1 = client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc("messages")
	ref = ref1.Collection("rawConversation").Doc("Context promt")
	_, err = ref.Set(ctx, map[string]interface{}{
		"content": `Eres un útil asistente de un restaurante diseñado para leer pedidos, compararlos con el menú "
		y generar el pedido en formato JSON, asegúrate de que cada platillo de la orden del cliente
		tenga los campos 'id', 'nombre_platillo', 'precio_por_cada_uno' y 'cantidad', debes devolver
		un JSON con el siguiente formato: ` + string(jsonOrdersData) + ` si el usuario no ordena nada,
		regresa el JSON vacío. Menu: ` + string(jsonMenuData) + `Se muy amigable, recuerda que nos puedes
		ayudar a conseguir mas clientes si les caes bien, y no pongas tanto texto, se amable pero conciso
		al mismo tiempo. Responde siempre en español y NO digas cosas como 'Aqui esta tu pedido en formato JSON'
		solo di 'Aqui esta tu pedido' o de alguna otra forma`,
		"role":      "system",
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		log.Fatalln(err)
	}
}
