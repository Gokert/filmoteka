package httpResponse

import (
	"encoding/json"
	"filmoteka/pkg/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func SendResponse(w http.ResponseWriter, r *http.Request, response *models.Response, log *logrus.Logger) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error("Send response error: ", err)
		response.Status = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Infof("Time: %v, Method: %s, Status: %d, URL: %s", time.Now(), r.Method, response.Status, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Error("Failed to send response: ", err.Error())
	}
}
