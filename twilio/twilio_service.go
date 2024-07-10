package twilioService

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func smsTwilio(message string) {
	from := os.Getenv("TWILIO_SMS_FROM")
	to := os.Getenv("TWILIO_SMS_TO")
	body := message
	var twilio_account_sid = os.Getenv("TWILIO_SID")
	var twilio_auth_token = os.Getenv("TWILIO_AUTH_TOKEN")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilio_account_sid,
		Password: twilio_auth_token,
	})

	params := &twilioApi.CreateMessageParams{
		To:   &to,
		From: &from,
		Body: &body,
	}

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending SMS message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}

}

func whatsAppTwilio(message string) {
	var twilio_account_sid = os.Getenv("TWILIO_SID")
	var twilio_auth_token = os.Getenv("TWILIO_AUTH_TOKEN")

	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilio_account_sid,
		Password: twilio_auth_token,
	})
	from := os.Getenv("TWILIO_WA_FROM")
	to := os.Getenv("TWILIO_WA_TO")
	body := message

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(body)

	_, err := twilioClient.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Message sent successfully!")
	}

}

// func ListenMsgs() {

// 	router := gin.Default()

// 	router.POST("/", func(context *gin.Context) {
// 		fileName, err := getNextAudioFileName()
// 		if err != nil {
// 			context.String(http.StatusInternalServerError, err.Error())
// 		}
// 		// Aquí es donde procesarías el mensaje recibido de WhatsApp
// 		// Puedes acceder al cuerpo del mensaje en context.Request.Body
// 		body, err := io.ReadAll(context.Request.Body)
// 		if err != nil {
// 			context.String(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		values, err := url.ParseQuery(string(body))
// 		if err != nil {
// 			context.String(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		var finalAnswer string
// 		// Extraer el valor del cuerpo del mensaje
// 		messageBody := values.Get("Body")
// 		contentType := values.Get("MediaContentType0")
// 		mediaUrl := values.Get("MediaUrl0")
// 		var textedQuestion string
// 		fmt.Println("contentType " + contentType)

// 		if contentType == "audio/ogg" {
// 			if err := downloadFile(mediaUrl, "audios/"+fileName); err != nil {
// 				fmt.Println("Error al descargar el archivo:", err)
// 				finalAnswer = "Error al descargar el archivo de audio"
// 			} else {
// 				textedQuestion = myOpenAi.Speech_to_text("audios/" + fileName)
// 			}
// 			fmt.Println("Media URL:  " + mediaUrl)
// 			textedQuestion = myOpenAi.Speech_to_text("audios/" + fileName)
// 			answerFromGPT := myOpenAi.AskGpt(textedQuestion)
// 			formatedAnswer := utils.FormatGPTResponse(answerFromGPT)

// 			finalAnswer = formatedAnswer

// 		} else {
// 			answerFromGPT := myOpenAi.AskGpt(messageBody)
// 			finalAnswer = utils.FormatGPTResponse(answerFromGPT)
// 			fmt.Println("MENSAJE: " + finalAnswer)
// 		}
// 		// Procesar el cuerpo del mensaje
// 		fmt.Println("Mensaje de WhatsApp recibido:", string(body))
// 		// Responder a la solicitud de ngrok
// 		context.String(http.StatusOK, finalAnswer)
// 	})

// 	// Iniciar el servidor en el puerto 3000
// 	if err := router.Run(":8000"); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func getNextAudioFileName() (string, error) {
// 	// Obtener la lista de archivos en el directorio 'audios'
// 	files, err := os.ReadDir("audios")
// 	if err != nil {
// 		return "", err
// 	}

// 	// Contar cuántos archivos de audio ya existen
// 	count := 0
// 	for _, file := range files {
// 		if strings.HasPrefix(file.Name(), "audio") && strings.HasSuffix(file.Name(), ".ogg") {
// 			count++
// 		}
// 	}

// 	// Crear el nombre para el próximo archivo de audio
// 	nextFileName := fmt.Sprintf("audio%d.ogg", count+1)
// 	return nextFileName, nil
// }

// func downloadFile(url string, filepath string) error {
// 	// Crear la carpeta 'audios' si no existe
// 	if err := os.MkdirAll("audios", os.ModePerm); err != nil {
// 		return err
// 	}

// 	// Crear el archivo para escribir el contenido del audio
// 	out, err := os.Create(filepath)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	// Crear la solicitud HTTP
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return err
// 	}

// 	// Agregar las credenciales de Twilio a la solicitud HTTP
// 	req.SetBasicAuth(os.Getenv("TWILIO_SID"), os.Getenv("TWILIO_AUTH_TOKEN"))

// 	// Realizar la solicitud HTTP
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	// Escribir el contenido del audio en el archivo
// 	_, err = io.Copy(out, resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }
