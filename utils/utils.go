package utils

import (
	firebaseService "baia_service/firebase"
	myOpenAi "baia_service/openai"
	"encoding/json"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
)

type Platillo struct {
	ID               int     `json:"id"`
	NombrePlatillo   string  `json:"nombre_platillo"`
	PrecioPorCadaUno float64 `json:"precio_por_cada_uno"`
	Cantidad         int     `json:"cantidad"`
}

type Orden struct {
	Orden []Platillo `json:"orden"`
}

func SendRequest(sentMessage string, senderID string, fbClient *firestore.Client) string {

	var finalAnswer string
	// Extraer el valor del cuerpo del mensaje
	firebaseService.SaveRawUserMessage(sentMessage, senderID, fbClient) // Use senderID from form values

	answerFromGPT := myOpenAi.AskGpt(sentMessage, senderID)
	firebaseService.SaveRawBAIAMessage(answerFromGPT, senderID, fbClient)

	finalAnswer = FormatGPTResponse(answerFromGPT)

	firebaseService.SaveBAIAMessage(finalAnswer, senderID, fbClient)

	firebaseService.SaveUserMessage(sentMessage, senderID, fbClient) // Use senderID from form values
	// Procesar el cuerpo del mensaje
	// Responder a la solicitud de ngrok
	return finalAnswer
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
	// JSON original

	// Parsear el JSON
	var orden Orden
	if err := json.Unmarshal([]byte(orderJson), &orden); err != nil {
		fmt.Println("ERROR AT PARSING JSON")
		return "", err
	}

	// Imprimir el detalle de cada platillo y calcular el total
	var output strings.Builder
	var total float64
	for _, platillo := range orden.Orden {
		subtotal := platillo.PrecioPorCadaUno * float64(platillo.Cantidad)
		total += subtotal
		output.WriteString(fmt.Sprintf("- %s (x%d): $%.2f\n", platillo.NombrePlatillo, platillo.Cantidad, subtotal))
	}

	output.WriteString(fmt.Sprintf("\nTotal del pedido: $%.2f\n", total))

	// Obtener la salida como una cadena
	result := output.String()
	return result, nil
}
