package firestoreService

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type UserData struct {
	NumberOfConversations int `firestore:"numberOfConvs"`
}

func createUserData(senderID string, client *firestore.Client) {
	ctx := context.Background()

	// Crear un nuevo documento con los datos iniciales
	_, err := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc("userData").Set(ctx, map[string]interface{}{
		"numberOfConvs": 1,
	})
	if err != nil {
		fmt.Printf("Error creating document: %v\n", err)
	}
}

func getNextConversationNumber(client *firestore.Client, senderID string) int {
	ctx := context.Background()
	dsnap, err := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc("userData").Get(ctx)
	if err != nil {
		fmt.Printf("Error getting document: %v\n", err)
		go createUserData(senderID, client)
		return 1
	}
	if !dsnap.Exists() {
		fmt.Println("Document does not exist!")
		go createUserData(senderID, client)
		return 1
	}

	fmt.Printf("Document data: %#v\n", dsnap.Data())

	var userData UserData
	err = dsnap.DataTo(&userData)
	if err != nil {
		fmt.Printf("Error mapping document data: %v\n", err)
		return 1
	}

	nextConversationNumber := userData.NumberOfConversations + 1

	_, err = client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc("userData").Update(ctx, []firestore.Update{
		{
			Path:  "numberOfConvs",
			Value: nextConversationNumber,
		},
	})
	if err != nil {
		fmt.Printf("Error updating document: %v\n", err)
		return userData.NumberOfConversations
	}

	return nextConversationNumber
}

func InitFirebase() *firestore.Client {
	ctx := context.Background()
	// config := &firebase.Config{
	// 	ProjectID: "baia-1df5a",
	// }
	sa := option.WithCredentialsFile("/Users/Armando09/Downloads/baia-1df5a-firebase-adminsdk-19h0s-8c00250774.json")

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Firestore initialized")
	return client
}

func getNextMessageNumber(client *firestore.Client, senderID string, isFromUser bool, isRaw bool, convoNumber int) int {
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
	collectionRef := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc(fmt.Sprintf("%v messages", convoNumber)).Collection(isRawMessage)
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

func SaveRawBAIAMessage(message string, senderID string, client *firestore.Client, convoNumber int) {
	ctx := context.Background()
	messageNumber := getNextMessageNumber(client, senderID, false, true, convoNumber)
	ref1 := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc(fmt.Sprintf("%v messages", convoNumber))
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

func SaveRawUserMessage(message string, senderID string, client *firestore.Client, convoNumber int) {
	ctx := context.Background()
	messageNumber := getNextMessageNumber(client, senderID, true, true, convoNumber)
	ref1 := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc(fmt.Sprintf("%v messages", convoNumber))
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

func SaveBAIAMessage(message string, senderID string, client *firestore.Client, convoNumber int) {
	ctx := context.Background()
	messageNumber := getNextMessageNumber(client, senderID, false, false, convoNumber)
	ref1 := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc(fmt.Sprintf("%v messages", convoNumber))
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

func SaveUserMessage(message string, senderID string, client *firestore.Client, convoNumber int) {
	ctx := context.Background()

	messageNumber := getNextMessageNumber(client, senderID, true, false, convoNumber)
	ref1 := client.Collection("El Sabor de Tlaxcala").Doc("conversations").Collection(senderID).Doc(fmt.Sprintf("%v messages", convoNumber))
	ref := ref1.Collection("Conversation").Doc(fmt.Sprintf("#%v message", messageNumber))
	_, err := ref.Set(ctx, map[string]interface{}{
		"content":   message,
		"role":      "user",
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("SAVED CLEAN USER")

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
		al mismo tiempo. Responde siempre en español y NO digas cosas como 'Aqui esta tu pedido en formato JSON'
		solo di 'Aqui esta tu pedido' o de alguna otra forma. Uno de tus propositos es incitar a comprar aunque te digan que no quieren.
		Cuando el usuario termine de hacer su orden, es muy importante que mandes un mensaje que diga "ORDEN COMPLETA", pero antes de eso tienes que
		preguntarle y asegurarte de que en serio haya acabado su pedido`,
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
		solo di 'Aqui esta tu pedido' o de alguna otra forma. Uno de tus propositos es incitar a comprar aunque te digan que no quieren.
		Cuando el usuario termine de hacer su orden, es muy importante que mandes un mensaje que diga "ORDEN COMPLETA", pero antes de eso tienes que
		preguntarle y asegurarte de que en serio haya acabado su pedido`,
		"role":      "system",
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		log.Fatalln(err)
	}
}
