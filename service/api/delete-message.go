package api

import (
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unsendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Autenticazione (Chi sei?)
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	myID, err := strconv.Atoi(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 2. Parsare l'ID del messaggio dall'URL
	// L'URL è /conversations/:id/messages/:messageId
	// Quindi usiamo ps.ByName("messageId")
	messageIDString := ps.ByName("messageId")
	messageID, err := strconv.Atoi(messageIDString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3. Chiamare il Database
	err = rt.db.DeleteMessage(messageID, myID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't delete message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Risposta successo (204 No Content)
	w.WriteHeader(http.StatusNoContent)
}