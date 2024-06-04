package utils

import (
	myOpenAi "baia_service/openai"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"firebase.google.com/go/v4/db"
)

type Message struct {
	Response   string `json:"response"`
	AfterOrder bool   `json:"afterOrder,omitempty"`
	IsImage    bool   `json:"isImage"`
}

type Order struct {
	ID               int     `json:"id"`
	NombrePlatillo   string  `json:"nombre_platillo"`
	PrecioPorCadaUno float64 `json:"precio_por_cada_uno"`
	Cantidad         int     `json:"cantidad"`
}

type Input struct {
	Messages []Message `json:"messages"`
	Orden    []Order   `json:"orden"`
}

type OutputMessage struct {
	Response string `json:"response"`
	IsImage  bool   `json:"isImage"`
}

type Output struct {
	Messages []OutputMessage `json:"messages"`
}

type Platillo struct {
	ID               int     `json:"id"`
	NombrePlatillo   string  `json:"nombre_platillo"`
	PrecioPorCadaUno float64 `json:"precio_por_cada_uno"`
	Cantidad         int     `json:"cantidad"`
}

type Orden struct {
	Orden []Platillo `json:"orden"`
}

func SendRequest(sentMessage string, senderID string, fbClient *db.Client) string {
	var finalAnswer Output

	// go realtimeService.SaveRawUserMessage(sentMessage, senderID, fbClient)

	answerFromGPT := myOpenAi.AskGpt(sentMessage, senderID)

	// go realtimeService.SaveRawBAIAMessage(answerFromGPT, senderID, fbClient)

	var input Input
	if err := json.Unmarshal([]byte(answerFromGPT), &input); err != nil {
		log.Printf("Error unmarshalling GPT response: %v", err)
		return answerFromGPT
	}

	finalAnswer = transform(input)

	// go realtimeService.SaveBAIAMessage(finalAnswer, senderID, fbClient)

	// go realtimeService.SaveUserMessage(sentMessage, senderID, fbClient) // Use senderID from form values

	finalAnswerJSON, err := json.Marshal(finalAnswer)
	if err != nil {
		log.Printf("Error marshalling final answer: %v", err)
		return answerFromGPT
	}

	return string(finalAnswerJSON)
}

func transform(input Input) Output {
	var output Output
	var orderSummary string
	var total float64

	for _, msg := range input.Messages {
		if !msg.AfterOrder {
			output.Messages = append(output.Messages, OutputMessage{
				Response: msg.Response,
				IsImage:  msg.IsImage,
			})
		}
	}

	for _, order := range input.Orden {
		orderSummary += fmt.Sprintf("- %s (x%d): $%.2f \n", order.NombrePlatillo, order.Cantidad, order.PrecioPorCadaUno*float64(order.Cantidad))
		total += order.PrecioPorCadaUno * float64(order.Cantidad)
	}

	if orderSummary != "" {
		output.Messages = append(output.Messages, OutputMessage{
			Response: orderSummary[:len(orderSummary)-1],
			IsImage:  false,
		})
	}

	for _, msg := range input.Messages {
		if msg.AfterOrder {
			if msg.Response == fmt.Sprintf("El total a pagar es: $%.1f", total) {
				output.Messages = append(output.Messages, OutputMessage{
					Response: msg.Response,
					IsImage:  msg.IsImage,
				})
			} else {
				output.Messages = append(output.Messages, OutputMessage{
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
	var orden Orden
	if err := json.Unmarshal([]byte(orderJson), &orden); err != nil {
		fmt.Println("ERROR AT PARSING JSON")
		return "", err
	}

	var output strings.Builder
	var total float64
	for _, platillo := range orden.Orden {
		subtotal := platillo.PrecioPorCadaUno * float64(platillo.Cantidad)
		total += subtotal
		output.WriteString(fmt.Sprintf("- %s (x%d): $%.2f\n", platillo.NombrePlatillo, platillo.Cantidad, subtotal))
	}

	output.WriteString(fmt.Sprintf("\nTotal del pedido: $%.2f\n", total))

	result := output.String()
	return result, nil
}
