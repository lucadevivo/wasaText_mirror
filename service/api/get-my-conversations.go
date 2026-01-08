package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Autenticazione
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

	// 2. Chiamata al DB
	conversations, err := rt.db.GetMyConversations(myID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Can't load conversations list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Risposta
	w.Header().Set("Content-Type", "application/json")
	if conversations == nil {
		conversations = []database.Conversation{}
	}
	_ = json.NewEncoder(w).Encode(conversations)
}