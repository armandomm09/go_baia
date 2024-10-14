package mongoService

import (
	baiaStructs "baia_service/structs"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserChatHistory(userID string, client *mongo.Client, Req openai.ChatCompletionRequest) openai.ChatCompletionRequest {
	ctx := context.Background()

	coll := client.Database("Sushi_Restaurant").Collection("Conversations")

	var conversation baiaStructs.Conversation
	err := coll.FindOne(ctx, bson.D{{"userID", userID}, {"isActive", true}}).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No active conversations found for user:", userID)
		} else {
			fmt.Println("Error getting the document:", err)
		}
		return Req
	}

	for _, message := range conversation.Messages {
		Req.Messages = append(Req.Messages, openai.ChatCompletionMessage{
			Role:    message.Role,
			Content: fmt.Sprintf("%v", message.Content),
		})
	}

	return Req
}

func SaveBAIAMessage(messages baiaStructs.FinalGPTResponse, senderID string, client *mongo.Client) {

	coll := client.Database("Sushi_Restaurant").Collection("Conversations")

	var activeConvResult bson.M

	err := coll.FindOne(context.TODO(), bson.D{{"isActive", true}, {"userID", senderID}}).
		Decode(&activeConvResult)
	if err == mongo.ErrNoDocuments {
		// fmt.Printf("No document was found")
		// return
		var conversation baiaStructs.Conversation
		conversation.UserID = senderID
		conversation.ID = uuid.New().String()
		conversation.IsActive = true

		var message baiaStructs.DBAssistantMessage
		message.Content = messages.Messages
		message.Role = "assistant"
		message.Timestamp = time.Now().Unix()

		conversation.Messages = append(conversation.Messages, message)

		result, err := coll.InsertOne(context.TODO(), conversation)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
		return
	}
	if err != nil {
		panic(err)
	}

	var messageToAdd baiaStructs.DBAssistantMessage

	messageToAdd.Content = messages.Messages
	messageToAdd.Role = "assistant"
	messageToAdd.Timestamp = time.Now().Unix()

	update := bson.D{
		{"$push", bson.D{
			{"messages", messageToAdd},
		}},
	}

	filter := bson.D{
		{"userID", senderID},
		{"isActive", true},
	}

	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message added to existing conversation")

}

func SaveUserMessage(message string, senderID string, client *mongo.Client) {

	coll := client.Database("Sushi_Restaurant").Collection("Conversations")

	var activeConvResult bson.M

	err := coll.FindOne(context.TODO(), bson.D{{"isActive", true}, {"userID", senderID}}).
		Decode(&activeConvResult)
	// if err == mongo.ErrNoDocuments {
	// 	// fmt.Printf("No document was found")
	// 	// return
	// 	var conversation baiaStructs.Conversation
	// 	conversation.UserID = senderID
	// 	conversation.ID = uuid.New().String()
	// 	conversation.IsActive = true

	// 	var message baiaStructs.DBAssistantMessage
	// 	message.Content = messages.Messages
	// 	message.Role = "assistant"
	// 	message.Timestamp = time.Now().Unix()

	// 	conversation.Messages = append(conversation.Messages, message)

	// 	result, err := coll.InsertOne(context.TODO(), conversation)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	// 	return
	// }
	// if err != nil {
	// 	panic(err)
	// }

	var messageToAdd baiaStructs.DBAssistantMessage

	// messageToAdd.Content = message
	var newOutputMessage baiaStructs.OutputMessage
	newOutputMessage.IsImage = false
	newOutputMessage.Response = message

	messageToAdd.Content = append(messageToAdd.Content, newOutputMessage)
	messageToAdd.Role = "user"
	messageToAdd.Timestamp = time.Now().Unix()

	update := bson.D{
		{"$push", bson.D{
			{"messages", messageToAdd},
		}},
	}

	filter := bson.D{
		{"userID", senderID},
		{"isActive", true},
	}

	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Message added to existing conversation")

}

func FinishOrder(serviceName string, userID string, order baiaStructs.Order, client *mongo.Client) {
	var finalOrder baiaStructs.FinalOrder

	finalOrder.ID = uuid.New().String()
	finalOrder.CreationDate = time.Now().Unix()
	finalOrder.State = "active"
	finalOrder.UserID = userID

	finalOrder.Order = order

	for i := 0; i < len(order.Order); i++ {
		finalOrder.Total += float32(order.Order[i].Quantity) * float32(order.Order[i].UnitaryPrice)
	}

	finalOrder.DeliveryLocation.Latitude = 19.041
	finalOrder.DeliveryLocation.Longitude = 98.206

	finalOrder.DeliveryAddress.Street = "Cda San Jose"
	finalOrder.DeliveryAddress.Number = 1411
	finalOrder.DeliveryAddress.Suburb = "Sta Cruz Buenavista"
	finalOrder.DeliveryAddress.Description = "Porton Gris"

	finalOrder.DeliveryDate = time.Now().Add(30 * time.Minute).Unix()

	finalOrder.Comments = "El sashimi que sea sin soya"

	coll := client.Database("Sushi_Restaurant").Collection("Orders")

	_, err := coll.InsertOne(context.TODO(), finalOrder)
	if err != nil {
		log.Fatal(err)
	}

	coll = client.Database("Sushi_Restaurant").Collection("Conversations")

	filter := bson.D{
		{"userID", userID},
		{"isActive", true},
	}

	// update := bson.D{
	// 	{"$set", bson.D{
	// 		{"isActive", false},
	// 	}},
	// }
	_, err = coll.InsertOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

}
