package realtimeService

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	db "firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

func InitFirebase() *db.Client {
	ctx := context.Background()

	// configure database URL
	conf := &firebase.Config{
		DatabaseURL: "https://baia-1df5a-default-rtdb.firebaseio.com/",
	}

	// fetch service account key
	opt := option.WithCredentialsFile("/Users/Armando09/Downloads/baia-1df5a-firebase-adminsdk-19h0s-8c00250774.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("error in initializing firebase app: ", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase DB client: ", err)
	}
	log.Println("Realtime Database Initialized")
	return client

}

func EncodeFirebaseKey(key string) string {
	key = strings.ReplaceAll(key, ".", "_")
	key = strings.ReplaceAll(key, "@", "_")
	key = strings.ReplaceAll(key, "/", "_")
	return key
}

func getNextMessageNumber(client *db.Client, senderID string, isFromUser bool, isRaw bool) (int, error) {
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

	ref := client.NewRef(fmt.Sprintf("El_Sabor_de_Tlaxcala/conversations/%s/messages/%s", senderID, isRawMessage))
	var messages map[string]interface{}
	if err := ref.Get(ctx, &messages); err != nil {
		return 0, err
	}

	messageCount := 0
	for key := range messages {
		if strings.Contains(key, messageType) {
			messageCount++
		}
	}
	return messageCount + 1, nil
}

func SaveRawBAIAMessage(message string, senderID string, client *db.Client) {
	ctx := context.Background()
	messageNumber, err := getNextMessageNumber(client, senderID, false, true)
	if err != nil {
		log.Println(err)
	}

	ref := client.NewRef(fmt.Sprintf("El_Sabor_de_Tlaxcala/conversations/%v/messages/rawConversation/%d response", senderID, messageNumber))
	if err := ref.Set(ctx, map[string]interface{}{
		"content":   message,
		"role":      "assistant",
		"timestamp": time.Now().Unix(),
	}); err != nil {
		log.Println(err)

	}

}

func SaveRawUserMessage(message string, senderID string, client *db.Client) {
	ctx := context.Background()
	// messageNumber, err := getNextMessageNumber(client, senderID, true, true)
	// // log.Printf("Raw user message sender ID: %v", senderID)
	// if err != nil {
	// 	return err
	// }

	refString := fmt.Sprintf("El_Sabor_de_Tlaxcala/conversations/%s/messages/rawConversation/%d message", senderID, 1)
	// log.Printf("REF STRING 1: \n%v", refString)
	ref := client.NewRef(refString)
	if err := ref.Set(ctx, map[string]interface{}{
		"content":   message,
		"role":      "user",
		"timestamp": time.Now().Unix(),
	}); err != nil {
		log.Println(err)

	}

}

func SaveBAIAMessage(message string, senderID string, client *db.Client) {
	ctx := context.Background()
	messageNumber, err := getNextMessageNumber(client, senderID, false, false)
	if err != nil {
		log.Println(err)
	}

	ref := client.NewRef(fmt.Sprintf("El_Sabor_de_Tlaxcala/conversations/%s/messages/Conversation/%d response", senderID, messageNumber))
	if err := ref.Set(ctx, map[string]interface{}{
		"content":   message,
		"role":      "assistant",
		"timestamp": time.Now().Unix(),
	}); err != nil {
		log.Println(err)

	}

}

func SaveUserMessage(message string, senderID string, client *db.Client) {
	ctx := context.Background()
	messageNumber, err := getNextMessageNumber(client, senderID, true, false)
	if err != nil {
		log.Println(err)
	}

	ref := client.NewRef(fmt.Sprintf("El_Sabor_de_Tlaxcala/conversations/%s/messages/Conversation/%d message", senderID, messageNumber))
	if err := ref.Set(ctx, map[string]interface{}{
		"content":   message,
		"role":      "user",
		"timestamp": time.Now().Unix(),
	}); err != nil {
		log.Println(err)
	}

}

func SetInitialPromt(senderID string, jsonMenuData []byte, jsonOrdersData []byte, client *db.Client) error {
	ctx := context.Background()

	prompt := `Eres un útil asistente de un restaurante diseñado para leer pedidos, compararlos con el menú "
	y generar el pedido en formato JSON, asegúrate de que cada platillo de la orden del cliente
	tenga los campos 'id', 'nombre_platillo', 'precio_por_cada_uno' y 'cantidad', debes devolver
	un JSON con el siguiente formato: ` + string(jsonOrdersData) + ` si el usuario no ordena nada,
	regresa el JSON vacío. Menu: ` + string(jsonMenuData) + `Se muy amigable, recuerda que nos puedes
	ayudar a conseguir mas clientes si les caes bien, y no pongas tanto texto, se amable pero conciso
	al mismo tiempo. Responde siempre en español y NO digas cosas como 'Aqui esta tu pedido en formato JSON'
	solo di 'Aqui esta tu pedido' o de alguna otra forma. Uno de tus propositos es incitar a comprar aunque te digan que no quieren.
	Cuando el usuario termine de hacer su orden, es muy importante que mandes un mensaje que diga "ORDEN COMPLETA", pero antes de eso tienes que
	preguntarle y asegurarte de que en serio haya acabado su pedido`

	ref := client.NewRef(fmt.Sprintf("El_Sabor_de_Tlaxcala/conversations/%s/messages/Conversation/Context promt", senderID))
	if err := ref.Set(ctx, map[string]interface{}{
		"content":   prompt,
		"role":      "system",
		"timestamp": time.Now().Unix(),
	}); err != nil {
		return err
	}

	ref = client.NewRef(fmt.Sprintf("El_Sabor_de_Tlaxcala/conversations/%s/messages/rawConversation/Context promt", senderID))
	if err := ref.Set(ctx, map[string]interface{}{
		"content":   prompt,
		"role":      "system",
		"timestamp": time.Now().Unix(),
	}); err != nil {
		return err
	}

	return nil
}
