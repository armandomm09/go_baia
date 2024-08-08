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
			fmt.Println("No active conversation found for user:", userID)
		} else {
			fmt.Println("Error getting document:", err)
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

		var message baiaStructs.DBMessage
		message.Content = messages.Messages
		message.Role = "user"
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

	var messageToAdd baiaStructs.DBMessage

	messageToAdd.Content = messages.Messages
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
