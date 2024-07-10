package baiaAPI

import (
	myOpenAi "baia_service/openai"
	"baia_service/utils"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
	"github.com/danielgtaylor/huma/v2"
)

type GPTResponse struct {
	Body struct {
		Answer string `json:answer`
	}
}

type UploadResponse struct {
	Body struct {
		Message string `json:"message"`
	}
}

type uploadFileRequest struct {
	Body struct {
		FileName string `json:"filename"`
	}
	RawBody multipart.Form
}

func RegisterEndPoints(api huma.API, fbClient *firestore.Client) {

	huma.Register(api, huma.Operation{
		OperationID:   "ask-about-order",
		Method:        http.MethodPost,
		Path:          "/baia/askGPT/text/{question}",
		Summary:       "Answers about your order",
		Tags:          []string{"BAIA"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *struct {
		Body *struct {
			Question string `json:"question" example:"Hola"`
			User     string `json:"senderID" example:"5212223201384@c.us"`
		}
	}) (*GPTResponse, error) {

		response := GPTResponse{}
		answer := utils.SendRequest(input.Body.Question, input.Body.User, fbClient)
		response.Body.Answer = answer

		fmt.Println("********** MESSAGE **********")
		fmt.Println(input.Body.User)
		fmt.Println(input.Body.Question)
		return &response, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "audio",
		Method:        http.MethodPost,
		Path:          "/baia/askGPT/audio/",
		Summary:       "Answers about your order sending the audio file",
		Tags:          []string{"BAIAudio"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *struct {
		// No Body, expecting multipart.Form directly
		RawBody multipart.Form
	}) (*GPTResponse, error) {
		// Verificar input.RawBody
		if input.RawBody.File == nil {
			return nil, huma.NewError(http.StatusBadRequest, "Request raw body is nil or does not contain files")
		}

		// Extract senderID from the form values
		senderID := input.RawBody.Value["senderID"]
		if len(senderID) == 0 {
			return nil, huma.NewError(http.StatusBadRequest, "Sender ID is missing")
		}

		if err := os.MkdirAll("audios", os.ModePerm); err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Error creating 'audios' directory", err)
		}

		fileHeaders, ok := input.RawBody.File["audio"]
		if !ok || len(fileHeaders) == 0 {
			return nil, huma.NewError(http.StatusBadRequest, "No audio file uploaded")
		}

		file, err := fileHeaders[0].Open()
		if err != nil {
			return nil, huma.NewError(http.StatusBadRequest, "Error opening uploaded file", err)
		}
		defer file.Close()

		dst, err := os.Create(filepath.Join("audios/apiAudios", fileHeaders[0].Filename))
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Error creating file on server", err)
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Error saving file to server", err)
		}

		audioPath := "audios/apiAudios/" + fileHeaders[0].Filename
		translatedText := myOpenAi.Speech_to_text(audioPath)

		// Verificar fbClient
		if fbClient == nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Firebase client is nil")
		}

		formatedAnswer := utils.SendRequest(translatedText, senderID[0], fbClient)

		response := GPTResponse{}
		response.Body.Answer = formatedAnswer

		return &response, nil
	})

}
