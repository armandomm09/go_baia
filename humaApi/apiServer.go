package baiaAPI

// import (
// 	"baia_service/firebase/firestoreService"
// 	"baia_service/firebase/realtimeService"
// 	myOpenAi "baia_service/openai"
// 	"baia_service/utils"
// 	"context"
// 	"encoding/json"
// 	"io"
// 	"log"
// 	"mime/multipart"
// 	"net/http"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	"cloud.google.com/go/firestore"
// 	"firebase.google.com/go/v4/db"
// 	"github.com/danielgtaylor/huma/v2"
// )

// type GPTResponse struct {
// 	Body struct {
// 		Answer string `json:answer`
// 	}
// }

// type AutoGenerated struct {
// 	Messages []struct {
// 		Response string `json:"response"`
// 		IsImage  bool   `json:"isImage"`
// 	} `json:"messages"`
// }

// type FinallOutput struct {
// 	Body struct {
// 		AutoGenerated `json:output`
// 	}
// }

// type UploadResponse struct {
// 	Body struct {
// 		Message string `json:"message"`
// 	}
// }

// type uploadFileRequest struct {
// 	Body struct {
// 		FileName string `json:"filename"`
// 	}
// 	RawBody multipart.Form
// }

// func RegisterEndPoints(api huma.API, rtClient *db.Client, fbClient *firestore.Client) {

// 	huma.Register(api, huma.Operation{
// 		OperationID:   "ask-about-order",
// 		Method:        http.MethodPost,
// 		Path:          "/baia/askGPT/text/{question}",
// 		Summary:       "Answers about your order",
// 		Tags:          []string{"BAIA"},
// 		DefaultStatus: http.StatusCreated,
// 	}, func(ctx context.Context, input *struct {
// 		Body *struct {
// 			Question string `json:"question" example:"Hola"`
// 			User     string `json:"senderID" example:"5212223201384@c.us"`
// 		}
// 	}) (*FinallOutput, error) {
// 		answer := utils.SendRequest(input.Body.Question, realtimeService.EncodeFirebaseKey(input.Body.User), rtClient)

// 		var gptResponse AutoGenerated
// 		var response FinallOutput
// 		if err := json.Unmarshal([]byte(answer), &gptResponse); err != nil {
// 			log.Printf("Error unmarshalling transformed response: %v", err)
// 			return nil, err
// 		}

// 		if strings.Contains(answer, "ORDEN COMPLETA") {
// 			go firestoreService.SaveOrderOnFirestore(realtimeService.EncodeFirebaseKey(input.Body.User), rtClient, fbClient)
// 		}

// 		log.Println("FINALL RESPONSE")
// 		log.Println(gptResponse)

// 		response.Body.AutoGenerated = gptResponse
// 		return &response, nil
// 	})

// 	huma.Register(api, huma.Operation{
// 		OperationID:   "audio",
// 		Method:        http.MethodPost,
// 		Path:          "/baia/askGPT/audio/",
// 		Summary:       "Answers about your order sending the audio file",
// 		Tags:          []string{"BAIAudio"},
// 		DefaultStatus: http.StatusCreated,
// 	}, func(ctx context.Context, input *struct {
// 		RawBody multipart.Form
// 	}) (*GPTResponse, error) {

// 		if input.RawBody.File == nil {
// 			return nil, huma.NewError(http.StatusBadRequest, "Request raw body is nil or does not contain files")
// 		}

// 		senderID := input.RawBody.Value["senderID"]
// 		if len(senderID) == 0 {
// 			return nil, huma.NewError(http.StatusBadRequest, "Sender ID is missing")
// 		}

// 		if err := os.MkdirAll("audios", os.ModePerm); err != nil {
// 			return nil, huma.NewError(http.StatusInternalServerError, "Error creating 'audios' directory", err)
// 		}

// 		fileHeaders, ok := input.RawBody.File["audio"]
// 		if !ok || len(fileHeaders) == 0 {
// 			return nil, huma.NewError(http.StatusBadRequest, "No audio file uploaded")
// 		}

// 		file, err := fileHeaders[0].Open()
// 		if err != nil {
// 			return nil, huma.NewError(http.StatusBadRequest, "Error opening uploaded file", err)
// 		}
// 		defer file.Close()

// 		dst, err := os.Create(filepath.Join("audios/apiAudios", fileHeaders[0].Filename))
// 		if err != nil {
// 			return nil, huma.NewError(http.StatusInternalServerError, "Error creating file on server", err)
// 		}
// 		defer dst.Close()

// 		_, err = io.Copy(dst, file)
// 		if err != nil {
// 			return nil, huma.NewError(http.StatusInternalServerError, "Error saving file to server", err)
// 		}

// 		audioPath := "audios/apiAudios/" + fileHeaders[0].Filename
// 		translatedText := myOpenAi.Speech_to_text(audioPath)

// 		if rtClient == nil {
// 			return nil, huma.NewError(http.StatusInternalServerError, "Firebase client is nil")
// 		}

// 		formatedAnswer := utils.SendRequest(translatedText, senderID[0], rtClient)

// 		response := GPTResponse{}
// 		response.Body.Answer = formatedAnswer

// 		return &response, nil
// 	})

// }