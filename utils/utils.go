package utils

import (
	"baia_service/mongoService"
	myOpenAi "baia_service/openai"
	baiaStructs "baia_service/structs"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

func SendRequest(sentMessage string, senderID string, client *mongo.Client) baiaStructs.FinalGPTResponse {
	var finalAnswer baiaStructs.FinalGPTResponse
	// go realtimeService.SaveRawUserMessage(sentMessage, senderID, fbClient)

	answerFromGPT := myOpenAi.AskGpt(sentMessage, senderID, client)

	// go realtimeService.SaveRawBAIAMessage(answerFromGPT, senderID, fbClient)

	var input baiaStructs.GPTUnformattedResponse
	if err := json.Unmarshal([]byte(answerFromGPT), &input); err != nil {
		log.Printf("Error unmarshalling GPT response: %v", err)
		var output baiaStructs.FinalGPTResponse
		output.Messages[0].Response = answerFromGPT
		return output
	}

	finalAnswer = transform(input)
	go mongoService.SaveBAIAMessage(finalAnswer, senderID, client)
	// go realtimeService.SaveBAIAMessage(finalAnswer, senderID, fbClient)

	// go realtimeService.SaveUserMessage(sentMessage, senderID, fbClient) // Use senderID from form values

	// finalAnswerJSON, err := json.Marshal(finalAnswer)
	// if err != nil {
	// 	log.Printf("Error marshalling final answer: %v", err)
	// 	return answerFromGPT
	// }

	return finalAnswer
}

func transform(input baiaStructs.GPTUnformattedResponse) baiaStructs.FinalGPTResponse {
	var output baiaStructs.FinalGPTResponse
	var orderSummary string
	var total float64

	for _, msg := range input.Messages {
		if !msg.AfterOrder {
			output.Messages = append(output.Messages, baiaStructs.OutputMessage{
				Response: msg.Response,
				IsImage:  msg.IsImage,
			})
		}
	}

	for _, order := range input.Order {
		orderSummary += fmt.Sprintf("- %s (x%d): $%.2f \n", order.NombrePlatillo, order.Cantidad, order.PrecioPorCadaUno*float64(order.Cantidad))
		total += order.PrecioPorCadaUno * float64(order.Cantidad)
	}

	if orderSummary != "" {
		output.Messages = append(output.Messages, baiaStructs.OutputMessage{
			Response: orderSummary[:len(orderSummary)-1],
			IsImage:  false,
		})
	}

	for _, msg := range input.Messages {
		if msg.AfterOrder {
			if msg.Response == fmt.Sprintf("El total a pagar es: $%.1f", total) {
				output.Messages = append(output.Messages, baiaStructs.OutputMessage{
					Response: msg.Response,
					IsImage:  msg.IsImage,
				})
			} else {
				output.Messages = append(output.Messages, baiaStructs.OutputMessage{
					Response: msg.Response,
					IsImage:  msg.IsImage,
				})
			}
		}
	}

	return output
}

func FormatGPTResponse(text string) string {
	if strings.Contains(text, "json") {
		jsonSubstring := strings.Split(text, "json")

		if strings.Contains(jsonSubstring[1], "]") {
			jsonSubstring2 := strings.Split(jsonSubstring[1], "]")
			pureJson := jsonSubstring2[0] + "]}"

			formatedJson, err := formatOrderFromJson(pureJson)
			if err != nil {
				return text
			}
			formatedText := jsonSubstring[0] + "\n \n" + "> " + formatedJson + "\n \n" + strings.TrimSpace(strings.Replace(jsonSubstring2[1], "}", "", -1))
			return strings.Replace(formatedText, "`", "", -1)
		}
		return text
	}

	return text
}

func formatOrderFromJson(orderJson string) (string, error) {
	var orden baiaStructs.Order
	if err := json.Unmarshal([]byte(orderJson), &orden); err != nil {
		fmt.Println("ERROR AT PARSING JSON")
		return "", err
	}

	var output strings.Builder
	var total float64
	for _, platillo := range orden.Order {
		subtotal := platillo.PrecioPorCadaUno * float64(platillo.Cantidad)
		total += subtotal
		output.WriteString(fmt.Sprintf("- %s (x%d): $%.2f\n", platillo.NombrePlatillo, platillo.Cantidad, subtotal))
	}

	output.WriteString(fmt.Sprintf("\nTotal del pedido: $%.2f\n", total))

	result := output.String()
	return result, nil
}
