package httpResponse

import (
	"encoding/json"
	"filmoteka/pkg/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func SendResponse(w http.ResponseWriter, r *http.Request, response *models.Response, log *logrus.Logger) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error("Send response error: ", err)
		response.Status = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Error("Failed to send response: ", err.Error())
	}
}
