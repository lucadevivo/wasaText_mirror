package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type MessageRequest struct {
	Content string `json:"content"`
	PhotoID int    `json:"photo_id"` // Nuovo campo opzionale
}

type MessageResponse struct {
	ID int `json:"id"`
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1. AUTENTICAZIONE DEL SENDER
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	senderID, err := strconv.Atoi(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 2. RECUPERO DESTINATARIO DALL'URL
	recipientIDString := ps.ByName("id")
	recipientID, err := strconv.Atoi(recipientIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3. PARSING DEL BODY
	var body MessageRequest
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validazione: Il messaggio deve avere o del testo o una foto.
	// Se entrambi mancano, è una richiesta non valida.
	if len(body.Content) == 0 && body.PhotoID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 4. ESECUZIONE (Database)
	// Passiamo senderID, recipientID, Content e anche PhotoID
	msgID, err := rt.db.SendMessage(senderID, recipientID, body.Content, body.PhotoID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error sending message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 5. RISPOSTA
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(MessageResponse{ID: msgID})
}
